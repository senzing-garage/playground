#!/usr/bin/env python3

# Import Python packages.

import grpc
from senzing_grpc import SzAbstractFactoryGrpc

grpc_channel = grpc.insecure_channel("localhost:8261")

# Create an abstract factory for accessing Senzing via gRPC.

sz_abstract_factory = SzAbstractFactoryGrpc(grpc_channel)

# Create Senzing object.

sz_engine = sz_abstract_factory.create_engine()

# List all methods for a Senzing object.

print(sz_engine.help())

# Print help for a specific method.

print(sz_engine.help("get_entity_by_record_id"))
