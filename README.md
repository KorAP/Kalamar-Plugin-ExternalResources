# Kalamar-Plugin-ExternalProvider

## Description

Kalamar-Plugin-ExternalProvider is a web service that integrates in the plugin framework of
[Kalamar](https://github.com/KorAP/Kalamar), to allow linking texts by their text sigle
to external data providers, mainly for full text access.

Kalamar-Plugin-Export is meant to be a basic export plugin and should
demonstrate and evaluate the plugin capabilities of Kalamar.

## Build

To build the latest version of Kalamar-Plugin-ExternalProvider, do ...

```shell
$ git clone https://github.com/KorAP/Kalamar-Plugin-ExternalProvider
$ cd Kalamar-Plugin-ExternalProvider
$ go build .
```

## Running

The binary can be started without prerequisites. The `templates` folder has to be kept in the root directory.

Registration of the plugin in Kalamar is not yet officially supported.
Registration works by passing the following JSON blob
to the plugin registration handler.

```json
{
  "name" : "External Provider",
  "desc" : "Buy content from an external provider",
  "embed" : [{
    "panel" : "match",
    "title" : "Full Text",
    "classes" : ["cart"],
    "onClick" : {
      "action" : "addWidget",
      "template" : "{SERVICE_URL}",
      "permissions": [
        "scripts",
        "popups" 
      ]
    }
  }]
}
```

## Customization

The following environment variables can be set either as environment variables
or via `.env` file.

- `KORAP_SERVER`: The server URL of the hosting service.
- `PORT`: The port the service should be listen to.
- `KORAP_EXTERNAL_PROVIDER`: The URL the service is hosted.

## License

Copyright (c) 2023, [IDS Mannheim](https://www.ids-mannheim.de/), Germany

Kalamar-Plugin-ExternalProvider is developed as part of the
[KorAP](https://korap.ids-mannheim.de/) Corpus Analysis Platform
at the Leibniz Institute for the German Language
([IDS](https://www.ids-mannheim.de/)).

Kalamar-Plugin-ExternalProvider is published under the
[BSD-2 License](https://raw.githubusercontent.com/KorAP/Kalamar-Plugin-Export/master/LICENSE).
