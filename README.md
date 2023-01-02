# Kalamar-Plugin-ExternalProvider

## Description

Kalamar-Plugin-ExternalProvider is a web service that integrates in the plugin framework of
[Kalamar](https://github.com/KorAP/Kalamar), to allow linking texts by their text sigle
to external data providers, mainly for full text access.

Kalamar-Plugin-ExternalProvider is meant to be a basic plugin and should
demonstrate and evaluate the plugin capabilities of Kalamar.

## Prerequisites

Go 1.19 or later

## Build

To build the latest version of Kalamar-Plugin-ExternalProvider, do ...

```shell
$ git clone https://github.com/KorAP/Kalamar-Plugin-ExternalProvider
$ cd Kalamar-Plugin-ExternalProvider
$ go test .
$ go build .
```

## Running

The binary can be started without prerequisites. The `templates` folder has to be kept in the root directory.

Registration of the plugin in Kalamar is not yet officially supported.
Registration works by passing the JSON blob generated at `/plugin.json`
to the plugin registration handler.

## Customization

The following environment variables can be set either as environment variables
or via `.env` file.

- `KORAP_SERVER`: The server URL of the hosting service.
- `KORAP_EXTERNAL_PROVIDER_PORT`: The port the service should be listen to.
- `KORAP_EXTERNAL_PROVIDER`: The exposed URL the service is hosted.

## License

Copyright (c) 2023, [IDS Mannheim](https://www.ids-mannheim.de/), Germany<br>
Author: [Nils Diewald](https://www.nils-diewald.de/)

Kalamar-Plugin-ExternalProvider is developed as part of the
[KorAP](https://korap.ids-mannheim.de/) Corpus Analysis Platform
at the Leibniz Institute for the German Language
([IDS](https://www.ids-mannheim.de/)).

Kalamar-Plugin-ExternalProvider is published under the
[BSD-2 License](https://raw.githubusercontent.com/KorAP/Kalamar-Plugin-ExternalProvider/master/LICENSE).
