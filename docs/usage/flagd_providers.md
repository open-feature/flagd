# Flagd Providers

Flagd providers are used for interacting with the `flagd` service via the OpenFeature SDK, they act as the translation layer between the evaluation API and the flag management system in use (in this case `flagd`).
Documentation for each language specific provider can be found below:

| Language      | Provider |
| ----------- | ----------- |
| Go      | [Go Flagd Provider](https://github.com/open-feature/go-sdk-contrib/tree/main/providers/flagd)
| Java   | [Java Flagd Provider](https://github.com/open-feature/java-sdk-contrib/tree/main/providers/flagd)
| Javascript   | [Javascript Flagd Provider](https://github.com/open-feature/js-sdk-contrib/tree/main/libs/providers/flagd)
| PHP   | [PHP Flagd Provider](https://github.com/open-feature/php-sdk-contrib/tree/main/src/Flagd)
| Python   | Not currently available, [help by contributing here](https://github.com/open-feature/python-sdk-contrib)
| .NET   | [.NET Flagd Provider](https://github.com/open-feature/dotnet-sdk-contrib/tree/main/src/OpenFeature.Contrib.Providers.Flagd)
| Ruby  | Not currently available, [help by contributing here](https://github.com/open-feature/ruby-sdk-contrib)

Any (new or existing) `flagd` providers ought to follow [these guidelines](../other_resources/creating_providers.md).
