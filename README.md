# Workload Orchestrator

A lightweight distributed workload orchestration system built with Golang.<div>
*This is a personal project.*

## Goals
- Coordinate Docker workloads across multiple machine nodes.
- Operate a central controller and multiple worker agents using gRPC
- Resource-aware workload placement
- Maintain desired states using local reconciliation by agents.
- Synchronize states between controller and agents using revisions, acknowledgments, and reconnection recovery.
- Detect any failures and recover workloads.

## Stack
- Language: Go
- Containers: Docker
- Communication: gRPC
- Persistence: SQLite
- Platform: Linux

## Status
**Planning / Early development** - project organization & core implementations.

Features and documentation will be added with development progress