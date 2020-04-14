![status](https://img.shields.io/badge/status-alpha-red.svg)
[![actions atatus](https://github.com/DamianSkrzypczak/shift/workflows/Pipeline/badge.svg)](https://github.com/DamianSkrzypczak/shift/actions)
[![go report](https://goreportcard.com/badge/github.com/DamianSkrzypczak/shift)](https://goreportcard.com/report/github.com/DamianSkrzypczak/shift)
[![godoc](https://godoc.org/github.com/DamianSkrzypczak/shift?status.svg)](http://godoc.org/github.com/DamianSkrzypczak/shift)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c)](https://pkg.go.dev/github.com/DamianSkrzypczak/shift)
[![license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/DamianSkrzypczak/shift/blob/master/LICENSE)

<img align="right" height="159px" src="./media/logo.png">

<h1>Shift</h1>
<p>Domain driven API framework written in Golang</p>


> :warning: **Solution currently unusable and in active development.**

## Idea

Framework is based on idea of fully code-separated domains.
The main benefit of such an approach would be the opportunity
to work with a monolithic codebase in the early stages of the
project and, over time, to move gradually towards architecture
based on microservices. See [project workflow graph](./media/workflow.png).

## Inter-domain communication

Ultimately, the framework is intended to provide tools
for inter-domain communication which will allow
automatic switching from in-memory communication
towards network-based communication between services
located on separate instances.

## What's next

Please note that framework is in very early stage so some basic components and features are still not there.

These parts are next in line:

- project & code documentation
- logging
- unit tests (including \_examples/)
- e2e bdd-style tests (including \_examples/)
- benchmarks

More about future plans soon
on project documentation page (in progress)
but right now please see [todo.md](./todo.md).
