# go-deploy
[![Build Status](https://travis-ci.org/UKFast-Mobile/go-deploy.svg?branch=master)](https://travis-ci.org/UKFast-Mobile/go-deploy)
[![Coverage Status](https://coveralls.io/repos/github/UKFast-Mobile/go-deploy/badge.svg?branch=master)](https://coveralls.io/github/UKFast-Mobile/go-deploy?branch=master)
[![codebeat badge](https://codebeat.co/badges/7176d1f3-efeb-4341-8c29-0a4c2777bb97)](https://codebeat.co/projects/github-com-ukfast-mobile-go-deploy-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/UKFast-Mobile/go-deploy)](https://goreportcard.com/report/github.com/UKFast-Mobile/go-deploy)
[![GoDoc](https://godoc.org/github.com/UKFast-Mobile/go-deploy?status.svg)](https://godoc.org/github.com/UKFast-Mobile/go-deploy)

> Deployment CLI tool

## :hammer: Work in progress

## Installation

```
 go get gopkg.in/UKFast-Mobile/go-deploy.v0
```

## Commands

### Setup

> Provides convenince CLI to setup a `go-deploy.json` file configuration

```
go-deploy setup staging
```

### Prepare

> Prepares remote server in accordance to the `go-deploy.json` configuration file

```
go-deploy prepare staging
```

### Deploy (main)

> Deploys your application in accordance ot the `go-deploy.json` configuration file

```
go-deploy staging
```

### Help

> Shows usage information for the go-deploy CLI

```
go-deploy help
```