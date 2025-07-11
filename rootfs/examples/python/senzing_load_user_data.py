#!/usr/bin/env python3

# Import Python packages.

import json

import grpc
from senzing import SzEngineFlags, SzError
from senzing_grpc import SzAbstractFactoryGrpc

# Set environment specific variables.

home_path = "./"

# Modify the following. Identify uploaded file.

filepath = home_path + "senzing-example-data.json"

# Discover `DATA_SOURCE` values in records.

datasources = []
with open(filepath, "r", encoding="utf-8") as file:
    for line in file:
        line_as_dict = json.loads(line)
        datasource = line_as_dict.get("DATA_SOURCE")
        if datasource not in datasources:
            datasources.append(datasource)
print(f"Found the following DATA_SOURCE values in the data: {datasources}")

# Create an abstract factory for accessing Senzing via gRPC.

grpc_channel = grpc.insecure_channel("localhost:8261")
sz_abstract_factory = SzAbstractFactoryGrpc(grpc_channel)

# Create Senzing objects.

sz_configmanager = sz_abstract_factory.create_configmanager()
sz_diagnostic = sz_abstract_factory.create_diagnostic()
sz_engine = sz_abstract_factory.create_engine()

# Get current Senzing configuration.

config_id = sz_configmanager.get_default_config_id()
sz_config = sz_configmanager.create_config_from_config_id(config_id)

# Add DataSources to Senzing configuration.

for datasource in datasources:
    try:
        sz_config.register_data_source(datasource)
    except SzError as err:
        print(err)

# Persist new Senzing configuration.

new_json_config = sz_config.export()
new_config_id = sz_configmanager.register_config(
    new_json_config, "Add user datasources"
)
sz_configmanager.replace_default_config_id(config_id, new_config_id)

# With the change in Senzing configuration, Senzing objects need to be updated.

sz_abstract_factory.reinitialize(new_config_id)

# Call Senzing to add records.

with open(filepath, "r") as file:
    for line in file:
        try:
            line_as_dict = json.loads(line)
            info = sz_engine.add_record(
                line_as_dict.get("DATA_SOURCE"),
                line_as_dict.get("RECORD_ID"),
                line,
                SzEngineFlags.SZ_WITH_INFO,
            )
            print(info)
        except SzError as err:
            print(err)
