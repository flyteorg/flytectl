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

💻 Contribute to API Documentation
==================================

1. Set up a local cluster by running ``flytectl sandbox start --source=$(pwd)`` in the root directory
2. xxx