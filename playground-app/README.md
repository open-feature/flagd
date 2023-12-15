# flagd playground

The flagd playground is an application designed to test the behavior of flagd.
It allows users to define flags and experiment with various inputs to understand how flagd responds.
This tool is particularly useful for developers and testers who are working with flagd in their projects and need a simple and effective way to validate their flag definitions.

## Development

### Getting Started

To get started with the development of the flagd playground, you'll need to have Node.js installed on your machine.

1. Install [Node.js](https://nodejs.org/en/download/) version 18 or newer.
1. From the root of the project, run `make playground-dev`.
1. Open your browser and navigate to [http://localhost:5173/](http://localhost:5173/);

> [!NOTE]
> This page is mostly unstyled because it inherits the styles from Mkdocs Material.

### Add a new scenario

A new scenario can be added to the playground by following these steps:

1. Add a new scenario file during the ``./src/scenarios`` directory.
1. Export a constant that conforms to the `Scenario` type.
1. Include the scenario in the scenarios objects at `./src/scenarios/index.ts`.

> [!NOTE]
> Make sure to update the docs once you're ready. This does not happen automatically! Please see below for more information.

### Adding Playground to the Docs

Adding the playground app to the docs can be done by running the following command from the root of the project:

```bash
make make playground-publish
```

> [!NOTE]
> This will build the app and copy the output to the docs.