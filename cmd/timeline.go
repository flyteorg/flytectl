package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/wcharczuk/go-chart/drawing"

	"github.com/wcharczuk/go-chart"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/service"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/lyft/flytestdlib/logger"

	adminIdl "github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	coreIdl "github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/lyft/flyteidl/clients/go/admin"
	"github.com/spf13/cobra"
)

type timelineFlags struct {
	persistentFlags
	ExecutionName *string
	OutputPath    *string
}

func newTimelineCmd(flags persistentFlags) *cobra.Command {
	timelineFlags := timelineFlags{persistentFlags: flags}
	timelineCmd := &cobra.Command{
		Use:   "timeline",
		Short: "Visualize workflow execution timeline.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			c := admin.InitializeAdminClient(ctx, *admin.GetConfig(ctx))
			return visualizeTimeline(ctx, c, timelineFlags)
		},
	}

	timelineFlags.ExecutionName = timelineCmd.Flags().String("execution", "", "Specifies the name of the execution to visualize.")
	timelineFlags.OutputPath = timelineCmd.Flags().String("output-path", "timeline.png", "Specifies the output image path.")

	return timelineCmd
}

func getStartedAtTime(nodeExec *adminIdl.NodeExecution) time.Time {
	if startedAt := nodeExec.Closure.StartedAt; startedAt != nil {
		return time.Unix(startedAt.Seconds, int64(startedAt.Nanos))
	} else if createdAt := nodeExec.Closure.CreatedAt; createdAt != nil {
		return time.Unix(createdAt.Seconds, int64(createdAt.Nanos))
	} else {
		return time.Now()
	}
}

func getEndTime(startedAt time.Time, d *duration.Duration) time.Time {
	if d == nil {
		return startedAt
	}

	goDuration, err := ptypes.Duration(d)
	if err != nil {
		logger.Errorf(context.TODO(), "Failed to parse duration [%v]. Error: %v", d, err)
		return startedAt
	}

	return startedAt.Add(goDuration)
}

func visualizeTimeline(ctx context.Context, adminClient service.AdminServiceClient, flags timelineFlags) error {
	chartTasks := make([]chart.StackedBar, 0, 10)
	token := ""
	firstTime := time.Now().Add(time.Hour * 10)
	lastTime := time.Unix(0, 0)
	barStyle := chart.Style{
		FillColor:   drawing.ColorFromHex("c11313"),
		StrokeColor: drawing.ColorFromHex("c11313"),
		StrokeWidth: 0,
	}

	noShowBarStyle := chart.Style{
		FillColor:   drawing.ColorFromHex("ffffff"),
		StrokeColor: drawing.ColorFromHex("ffffff"),
		StrokeWidth: 0,
	}

	allResp := make([]*adminIdl.NodeExecution, 0, 100)

	for {
		resp, err := adminClient.ListNodeExecutions(ctx, &adminIdl.NodeExecutionListRequest{
			WorkflowExecutionId: &coreIdl.WorkflowExecutionIdentifier{
				Project: *flags.Project,
				Domain:  *flags.Domain,
				Name:    *flags.ExecutionName,
			},
			Limit: 100,
			Token: token,
		})

		if err != nil {
			return err
		}

		allResp = append(allResp, resp.NodeExecutions...)

		if len(resp.GetToken()) == 0 {
			break
		}
	}

	for _, nodeExec := range allResp {
		startedAt := getStartedAtTime(nodeExec)
		finishedAt := getEndTime(startedAt, nodeExec.Closure.Duration)
		if firstTime.After(startedAt) {
			firstTime = startedAt
		}

		if finishedAt.After(lastTime) {
			lastTime = finishedAt
		}
	}

	for i, nodeExec := range allResp {
		startedAt := getStartedAtTime(nodeExec)
		finishedAt := getEndTime(startedAt, nodeExec.Closure.Duration)
		chartTasks = append(chartTasks, chart.StackedBar{
			Name: strconv.Itoa(i),
			Values: []chart.Value{
				{
					Style: noShowBarStyle,
					Value: startedAt.Sub(firstTime).Minutes(),
				},
				{
					Style: barStyle,
					Value: finishedAt.Sub(startedAt).Minutes(),
				},
				{
					Style: noShowBarStyle,
					Value: lastTime.Sub(finishedAt).Minutes(),
				},
			},
		})
	}

	chartData := chart.StackedBarChart{
		Title:      fmt.Sprintf("%v-%v-%v", *flags.Project, *flags.Domain, *flags.ExecutionName),
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Width:  8096,
		Height: 1024,
		Bars:   chartTasks,
		XAxis:  chart.StyleShow(),
		YAxis:  chart.StyleShow(),
		//Height:     (chartTasks[0].GetHeight() + 30) * len(chartTasks),
		//BarSpacing: 30,
	}

	if flags.OutputPath != nil {
		f, err := os.Create(*flags.OutputPath)
		if err != nil {
			return err
		}

		defer func() {
			err = f.Close()
			if err != nil {
				panic(err)
			}
		}()

		return chartData.Render(chart.PNG, f)
	}

	return nil

	//result := render.ProcessStructured(firstTime, chartData)
	//logger.Print(ctx, result.Code)
	//
	//if result.Code > 0 {
	//	return fmt.Errorf(result.Message)
	//}
	//
	////save to file
	//if flags.OutputPath != nil {
	//	return result.Context.SavePNG(*flags.OutputPath)
	//}
	//
}
