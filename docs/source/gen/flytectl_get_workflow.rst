.. _flytectl_get_workflow:

flytectl get workflow
---------------------

Gets workflow resources

Synopsis
~~~~~~~~



Retrieves all the workflows within project and domain.(workflow,workflows can be used interchangeably in these commands)
::

 flytectl get workflow -p flytesnacks -d development

Retrieves workflow by name within project and domain.

::

 flytectl get workflow -p flytesnacks -d development  core.basic.lp.go_greet

Retrieves latest version of workflow by name within project and domain.

::

 flytectl get workflow -p flytesnacks -d development  core.basic.lp.go_greet --latest

Retrieves particular version of workflow by name within project and domain.

::

 flytectl get workflow -p flytesnacks -d development  core.basic.lp.go_greet --version v2

Retrieves all the workflows with filters.
::
 
  bin/flytectl get workflow -p flytesnacks -d development  --filter.field-selector="workflow.name=k8s_spark.dataframe_passing.my_smart_schema"
 
Retrieve specific workflow with filters.
::
 
  bin/flytectl get workflow -p flytesnacks -d development k8s_spark.dataframe_passing.my_smart_schema --filter.field-selector="workflow.version=v1"
  
Retrieves all the workflows with limit and sorting.
::
  
  bin/flytectl get -p flytesnacks -d development workflow  --filter.sort-by=created_at --filter.limit=1 --filter.asc

Retrieves all the workflow within project and domain in yaml format.

::

 flytectl get workflow -p flytesnacks -d development -o yaml

Retrieves all the workflow within project and domain in json format.

::

 flytectl get workflow -p flytesnacks -d development -o json

Usage


::

  flytectl get workflow [flags]

Options
~~~~~~~

::

      --filter.asc                     Specifies the sorting order. By default flytectl sort result in descending order
      --filter.field-selector string   Specifies the Field selector
      --filter.limit int32             Specifies the limit (default 100)
      --filter.sort-by string          Specifies which field to sort results  (default "created_at")
  -h, --help                           help for workflow

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --admin.authorizationHeader string           Custom metadata header to pass JWT
      --admin.authorizationServerUrl string        This is the URL to your IdP's authorization server. It'll default to Endpoint
      --admin.clientId string                      Client ID (default "flytepropeller")
      --admin.clientSecretLocation string          File containing the client secret (default "/etc/secrets/client_secret")
      --admin.endpoint string                      For admin types,  specify where the uri of the service is located.
      --admin.insecure                             Use insecure connection.
      --admin.maxBackoffDelay string               Max delay for grpc backoff (default "8s")
      --admin.maxRetries int                       Max number of gRPC retries (default 4)
      --admin.perRetryTimeout string               gRPC per retry timeout (default "15s")
      --admin.scopes strings                       List of scopes to request
      --admin.tokenUrl string                      OPTIONAL: Your IdP's token endpoint. It'll be discovered from flyte admin's OAuth Metadata endpoint if not provided.
      --admin.useAuth                              Deprecated: Auth will be enabled/disabled based on admin's dynamically discovered information.
      --config string                              config file (default is $HOME/.flyte/config.yaml)
  -d, --domain string                              Specifies the Flyte project's domain.
      --logger.formatter.type string               Sets logging format type. (default "json")
      --logger.level int                           Sets the minimum logging level. (default 4)
      --logger.mute                                Mutes all logs regardless of severity. Intended for benchmarks/tests only.
      --logger.show-source                         Includes source code location in logs.
  -o, --output string                              Specifies the output type - supported formats [TABLE JSON YAML] (default "TABLE")
  -p, --project string                             Specifies the Flyte project.
      --root.domain string                         Specified the domain to work on.
      --root.output string                         Specified the output type.
      --root.project string                        Specifies the project to work on.
      --storage.cache.max_size_mbs int             Maximum size of the cache where the Blob store data is cached in-memory. If not specified or set to 0,  cache is not used
      --storage.cache.target_gc_percent int        Sets the garbage collection target percentage.
      --storage.connection.access-key string       Access key to use. Only required when authtype is set to accesskey.
      --storage.connection.auth-type string        Auth Type to use [iam, accesskey]. (default "iam")
      --storage.connection.disable-ssl             Disables SSL connection. Should only be used for development.
      --storage.connection.endpoint string         URL for storage client to connect to.
      --storage.connection.region string           Region to connect to. (default "us-east-1")
      --storage.connection.secret-key string       Secret to use when accesskey is set.
      --storage.container string                   Initial container to create -if it doesn't exist-.'
      --storage.defaultHttpClient.timeout string   Sets time out on the http client. (default "0s")
      --storage.enable-multicontainer              If this is true,  then the container argument is overlooked and redundant. This config will automatically open new connections to new containers/buckets as they are encountered
      --storage.limits.maxDownloadMBs int          Maximum allowed download size (in MBs) per call. (default 2)
      --storage.type string                        Sets the type of storage to configure [s3/minio/local/mem/stow]. (default "s3")

SEE ALSO
~~~~~~~~

* :doc:`flytectl_get` 	 - Used for fetching various flyte resources including tasks/workflows/launchplans/executions/project.

