#!/usr/bin/env python3

# Load Senzing truth-sets

import json
import shutil

import grpc
import requests
from senzing import SzEngineFlags, SzError
from senzing_grpc import SzAbstractFactoryGrpc

# Set environment specific variables.

home_path = "./"
truth_set_url_prefix = "https://raw.githubusercontent.com/Senzing/truth-sets/refs/heads/main/truthsets/demo/"
truth_set_filenames = [
    "customers.json",
    "reference.json",
    "watchlist.json",
]

# Download truth-set files.

for filename in truth_set_filenames:
    url = truth_set_url_prefix + filename
    filepath = home_path + filename
    response = requests.get(url, stream=True, timeout=10)
    response.raw.decode_content = True
    with open(filepath, "wb") as file:
        shutil.copyfileobj(response.raw, file)

# Discover `DATA_SOURCE` values in records.

datasources = []
for filename in truth_set_filenames:
    filepath = home_path + filename
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
    new_json_config, "Add TruthSet datasources"
)
sz_configmanager.replace_default_config_id(config_id, new_config_id)

# With the change in Senzing configuration, Senzing objects need to be updated.

sz_abstract_factory.reinitialize(new_config_id)

# Call Senzing to add records.

for filename in truth_set_filenames:
    filepath = home_path + filename
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

# Retrieve an entity by identifying a record of the entity.

customer_1070_entity = sz_engine.get_entity_by_record_id("CUSTOMERS", "1070")
print(json.dumps(json.loads(customer_1070_entity), indent=2))

# Search for entities by attributes.

search_query = {
    "name_full": "robert smith",
    "date_of_birth": "11/12/1978",
}
search_result = sz_engine.search_by_attributes(json.dumps(search_query))
print(json.dumps(json.loads(search_result), indent=2))
