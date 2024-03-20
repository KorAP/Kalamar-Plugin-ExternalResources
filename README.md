# Kalamar-Plugin-ExternalResources

[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.8425526.svg)](https://doi.org/10.5281/zenodo.8425526)

## Description

`Kalamar-Plugin-ExternalResources` is a web service that integrates in the plugin framework of
[Kalamar](https://github.com/KorAP/Kalamar), to allow linking texts by their text sigle
to external data providers, mainly for full text access.

`Kalamar-Plugin-ExternalResources` is meant to be a basic plugin and should
demonstrate and evaluate the plugin capabilities of Kalamar.

**This is software in an early development stage. Its behaviour may change without warnings!**

## Prerequisites

Go 1.19 or later

## Build

To build the latest version of `Kalamar-Plugin-ExternalResources`, do ...

```shell
git clone https://github.com/KorAP/Kalamar-Plugin-ExternalResources
cd Kalamar-Plugin-ExternalResources
go test .
go build .
```

## Running

The binary can be started without prerequisites.
The `templates` folder has to be kept in the root directory.
A `db` folder contains the database.

Registration of the plugin in Kalamar is not yet officially supported -
but it works by passing the JSON blob generated at `/plugin.json`
to the plugin registration handler.

## Indexation

To index further data, the mappings need to be stored in a `csv`-file, like

```csv
"WPD11/A00/00001","Wikipedia","http://de.wikipedia.org/wiki/Alan_Smithee"
"WPD11/A00/00003","Wikipedia","http://de.wikipedia.org/wiki/Actinium"
"WPD11/A00/00005","Wikipedia","http://de.wikipedia.org/wiki/Ang_Lee"
```

With the first column being the textSigle (aka the document identifier),
the second being the provider name and the third being the URL.
These files can be gzipped as well.
Then run the indexation with:

```shell
./Kalamar-Plugin-ExternalResources data.csv
```

## Customization

The following environment variables can be set either as environment variables
or via `.env` file in the calling directory.

- `KORAP_SERVER`: The server URL of the hosting service.
- `KORAP_EXTERNAL_RESOURCES_PORT`: The port the service should be listen to.
- `KORAP_EXTERNAL_RESOURCES`: The exposed URL the service is hosted.

## Dockerization

Currently no official Docker image is provided.
To build an image based on the provided Dockerfile, run

```shell
docker build \
       -f Dockerfile \
       -t korap/kalamar-plugin-externalresources .
```

To create a container on Linux based on the image
with a mounted database in `db`
and a configuration file, run

```shell
docker run \
       --rm \
       --network host \
       -v ${PWD}/db/:/db/:z \
       -v ${PWD}/.env:/.env korap/ \
       kalamar-plugin-externalresources
```

To index using docker, run

```shell
docker run \
       --rm \
       --network host \
       -v ${PWD}/db/:/db/:z \
       -v ${PWD}/data/:/data \
       -v ${PWD}/.env:/.env \
       korap/kalamar-plugin-externalresources \
       data/data.csv.gz
```

## License

Copyright (c) 2023-2024, [IDS Mannheim](https://www.ids-mannheim.de/), Germany<br>
Author: [Nils Diewald](https://www.nils-diewald.de/)

Kalamar-Plugin-ExternalResources is developed as part of the
[KorAP](https://korap.ids-mannheim.de/) Corpus Analysis Platform
at the Leibniz Institute for the German Language
([IDS](https://www.ids-mannheim.de/)).

Kalamar-Plugin-ExternalResources is published under the
[BSD-2 License](https://raw.githubusercontent.com/KorAP/Kalamar-Plugin-ExternalResources/master/LICENSE).
