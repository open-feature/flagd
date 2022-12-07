# Docs

This directory contains all flagd documentation, see table of contents below:

## Usage

There are many ways to get started with flagd, the sections below run through some simple deployment options, followed by example evaluation calls made directly to the flagd instance via curl.

- Getting started
- Evaluation examples 

## Flag Configuration

Flagd is configured via CLI arguments on startup, these configuration options can be found in the Flagd Configuration section. The remaining sections cover the flag configurations themselves, which are JSON representations of the flag variants and targeting rules.

- Flagd Configuration
- Flag configuration 
- Fractional evaluation
- Reusable targeting rules

## Help

This section documents any behavior of flagd which may seem unexpected, currently covering 2 topics; why the HTTP int response is a string, and why values may be omitted from the evaluation response.

- HTTP int response 
- Omitted value from evaluation response

## Other Resources
- High level architecture 
- Caching 
- Snap
- Systemd service 