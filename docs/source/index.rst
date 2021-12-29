.. flytectl doc

######################
Flytectl Reference
######################

Overview
=========
This video will take you on a tour of Flytectl - how to install and configure it, as well as how to use the Verbs and Nouns sections on the left hand side menu. Detailed information can be found in the sections below the video.

.. youtube:: cV8ezYnBANE


Install
=======
Flytectl is a Golang binary and can be installed on any platform supported by
golang

.. tabbed:: OSX

  .. prompt:: bash $

      brew install flyteorg/homebrew-tap/flytectl

  *Upgrade* existing installation using the following command:

  .. prompt:: bash $

      brew update && brew upgrade flytectl

.. tabbed:: Other Operating systems

  .. prompt:: bash $

      curl -sL https://ctl.flyte.org/install | bash

**Test** if Flytectl is installed correctly (your Flytectl version should be > 0.2.0) using the following command:

.. prompt:: bash $

  flytectl version

Configure
=========
Flytectl allows configuring using a YAML file or pass every configuration value
on command-line. The following configuration is useful to setup.

Basic Configuration
--------------------

Flytectl configuration. 
For full list of configurable options available can be found by running ``flytectl --help`` and can be alternately be found `here <https://docs.flyte.org/projects/flytectl/en/stable/gen/flytectl.html#synopsis>`__

.. NOTE::

     Only Project(-p), Domain(-d), Output(-o) currently cannot be used in config file

.. tabbed:: Local Flyte Sandbox

    Automatically configured for you by ``flytectl sandbox`` command.

    .. code-block:: yaml

        admin:
          # For GRPC endpoints you might want to use dns:///flyte.myexample.com
          endpoint: dns:///localhost:30081
          insecure: true # Set to false to enable TLS/SSL connection (not recommended except on local sandbox deployment)
          authType: Pkce # authType: Pkce # if using authentication or just drop this.
        storage:
          connection:
            access-key: minio
            auth-type: accesskey
            disable-ssl: true
            endpoint: http://localhost:30084
            region: my-region-here
            secret-key: miniostorage
          container: my-s3-bucket
          type: minio

.. tabbed:: AWS Configuration

    .. code-block:: yaml

        admin:
          # For GRPC endpoints you might want to use dns:///flyte.myexample.com
          endpoint: dns:///<replace-me>
          authType: Pkce # authType: Pkce # if using authentication or just drop this.
          insecure: true # insecure: True # Set to true if the endpoint isn't accessible through TLS/SSL connection (not recommended except on local sandbox deployment)
        storage:
          type: stow
          stow:
            kind: s3
            config:
                auth_type: iam
                region: <REGION> # Example: us-east-2
          container: <replace> # Example my-bucket. Flyte k8s cluster / service account for execution should have read access to this bucket

.. tabbed:: GCS Configuration

    .. code-block:: yaml

        admin:
          # For GRPC endpoints you might want to use dns:///flyte.myexample.com
          endpoint: dns:///<replace-me>
          authType: Pkce # authType: Pkce # if using authentication or just drop this.
          insecure: false # insecure: True # Set to true if the endpoint isn't accessible through TLS/SSL connection (not recommended except on local sandbox deployment)
        storage:
          type: stow
          stow:
            kind: google
            config:
                json: ""
                project_id: <replace-me> # TODO: replace <project-id> with the GCP project ID
                scopes: https://www.googleapis.com/auth/devstorage.read_write
          container: <replace> # Example my-bucket. Flyte k8s cluster / service account for execution should have access to this bucket

.. tabbed:: Others

    For other supported storage backends like Oracle, Azure, etc., refer to the configuration structure `here <https://pkg.go.dev/github.com/flyteorg/flytestdlib/storage#Config>`__.

    Place the config file in ``$HOME/.flyte`` directory with the name config.yaml.
    This file is typically searched in:

    * ``$HOME/.flyte``
    * currDir from where you run flytectl
    * ``/etc/flyte/config``
    
    You can pass the file name in the commandline using ``--config <config-file-path>`` as well!


.. toctree::
   :maxdepth: 1
   :hidden:

   |plane| Getting Started <https://docs.flyte.org/en/latest/getting_started.html>
   |book-reader| User Guide <https://docs.flyte.org/projects/cookbook/en/latest/user_guide.html>
   |chalkboard| Tutorials <https://docs.flyte.org/projects/cookbook/en/latest/tutorials.html>
   |project-diagram| Concepts <https://docs.flyte.org/en/latest/concepts/basics.html>
   |rocket| Deployment <https://docs.flyte.org/en/latest/deployment/index.html>
   |book| API Reference <https://docs.flyte.org/en/latest/reference/index.html>
   |hands-helping| Community <https://docs.flyte.org/en/latest/community/index.html>

.. toctree::
   :maxdepth: -1
   :caption: Flytectl
   :hidden:

   Install and Configure <self>
   verbs
   nouns
   contribute
