# Toy KV

An attempt to build a toy key-value storage with durability.

Assumptions:
- Working set fits into memory
- There are no concurrent operations
- Writes are atomic

The goal is to build a durable key-value storage. Issues of isolation and failure recovery are ignored completely, at least for now.