.. flytectl doc

##########################################
Welcome to ``Flytectl``'s documentation!
##########################################


Installation
=============
Flytectl is a Golang binary and can be installed on any platform supported by
golang. To install simply copy paste the following into the command-line

.. prompt:: bash

   curl -s https://raw.githubusercontent.com/lyft/flytectl/master/install.sh | bash


Configuration
==============
Flytectl allows configuring using a YAML file or pass every configuration value
on command-line. The follow configuration is useful to setup.

Basic Configuration
--------------------

.. code-block:: yaml

  admin:
    # For GRPC endpoints you might want to use dns:///flyte.myexample.com
    endpoint: dns:///flyte.lyft.net
    # Change insecure flag to ensure that you use the right setting for your environment
    insecure: true
  # Logger settings to control logger output. Useful to debug
  #logger:
    #show-source: true
    #level: 1



.. toctree::
   :maxdepth: 1
   :caption: Flyte Core docs

   Flyte Documentation <https://flyte.readthedocs.io/en/latest/>

   gen/flytectl
.. toctree::
   :maxdepth: 1
   :caption: Flytectl verbs

   gen/flytectl_create
   gen/flytectl_get
   gen/flytectl_update
   gen/flytectl_delete
   gen/flytectl_register
   gen/flytectl_config

.. toctree::
   :maxdepth: 1
   :caption: Flytectl nouns

   gen/flytectl_create_project
   gen/flytectl_create_execution
   gen/flytectl_get_execution
   gen/flytectl_get_project
   gen/flytectl_get_workflow
   gen/flytectl_get_task
   gen/flytectl_get_launchplan
   gen/flytectl_update_project
   gen/flytectl_register_files
   gen/flytectl_version
   gen/flytectl_config_validate

.. toctree::
   :maxdepth: 2
   :caption: Contribute

   contribute


Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
