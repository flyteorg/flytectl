###########################
FlyteCTL Contribution Guide
###########################

First off, thank you for thinking about contributing! 
Below you’ll find instructions that will hopefully guide you through how to contribute to, fix, and improve FlyteCTL.

📝 Contribute to Documentation
==============================

Docs are generated using Sphinx and are available at [flytectl.rtfd.io](https://flytectl.rtfd.io).

To update the documentation, follow these steps:

1. Install the requirements by running ``pip install -r doc-requirements.txt`` in the root folder
2. Make modifications in the `docs/source/gen <https://github.com/flyteorg/flytectl/tree/master/docs/source/gen>`__ folder
3. Run ``make gendocs`` from within the `docs <https://github.com/flyteorg/flytectl/tree/master/docs>`__ folder
4. Open html files produced by Sphinx in your browser to verify if the changes look as expected (html files can be found in the ``docs/build/html`` folder)

💻 Contribute Code
==================

1. Set up a local cluster by running ``flytectl sandbox start --source=$(pwd)`` in the root directory
2. Run ``make compile`` in the root directory to compile the code
3. Run ``flytectl get project`` to see if things are working
4. Run the command you want to test in the terminal
5. If you want to update the command (add additional options, change existing options, etc.):
   * Navigate to `cmd <https://github.com/flyteorg/flytectl/tree/master/cmd>`__ directory
   * Each sub-directory points to a command, for example, ``create`` points to ``flytectl create ...``
   * Here are the directories you can navigate to:
     .. list-table:: Title
        :widths: 25 25 50
        :header-rows: 1

        * - Directory
          - Command
          - Description
        * - ``config``
          - ``flytectl config ...``
          - Common package for all commands; has root flags
        * - ``configuration``
          - ``flytectl configuration ...``
          - Command to validate/generate flytectl config
        * - ``create``
          - ``flytectl create ...``
          - Command to create a task/workflow/launchplan/execution/project
        * - ``delete``
          - ``flytectl delete ...``
          - Command to delete a task/workflow/launchplan/execution/project
        * - ``get``
          - ``flytectl get ...``
          - Command to get a task/workflow/launchplan/execution/project
        * - ``register``
          - ``flytectl register ...``
          - Command to register a task/workflow/launchplan
        * - ``sandbox``
          - ``flytectl sandbox ...``
          - Command to interact with sandbox
        * - ``update``
          - ``flytectl update ...``
          - Command to update a project
        * - ``upgrade``
          - ``flytectl upgrade ...``
          - Command to upgrade/rollback FlyteCTL version
        * - ``version``
          - ``flytectl version ...``
          - Command to fetch FlyteCTL version
   * Run appropriate tests to test the changes by running ``go test ./... -race -coverprofile=coverage.txt -covermode=atomic  -v`` 
     in the root directory