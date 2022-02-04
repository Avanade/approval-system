# github-mgmt
> This tool allows an organization to manage users, and their association to a public GitHub, and a private GitHub enterprise (for InnerSource.)

[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)
![GitHub issues](https://img.shields.io/github/issues/Avanade/avanade-template)
![GitHub](https://img.shields.io/github/license/Avanade/avanade-template)
![GitHub Repo stars](https://img.shields.io/github/stars/Avanade/avanade-template?style=social)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](https://avanade.github.io/code-of-conduct/)
[![Incubating InnerSource](https://img.shields.io/badge/Incubating-Ava--Maturity-%23FF5800?labelColor=yellow)](https://avanade.github.io/maturity-model/)

## Overview
This project contains a PowerApp, which is the main entry point for the application. Backend functionality is provided through a Go application, which handles associations with GitHub, and events (including leavers.)

It's recommended to call the leaving API as part of your JML process when users leave the organisation, but a Power Automate example is provided as an alternative (this is significantly less performant.)

## Licensing
github-mgmt is available under the [MIT Licence](LICENCE).

## Solutions Referenced

- [Azure SQL Database ledger tables](https://docs.microsoft.com/en-us/azure/azure-sql/database/ledger-overview?WT.mc_id=AI-MVP-5004204)
- [Azure Confidential Ledger](https://docs.microsoft.com/en-gb/azure/confidential-ledger/?WT.mc_id=AI-MVP-5004204)


```
These are provided as examples. Include links to components you have used, or delete this section.
DELETE THIS COMMENT
```

## Documentation
The `docs` folder contains [more detailed documentation](./docs/start-here.md), along with setup instructions.

## Contact
Feel free to [raise an issue on GitHub](https://github.com/Avanade/avanade-template/issues), or see our [security disclosure](SECURITY.md) policy.

## Contributing
Contributions are welcome. See information on [contributing](./CONTRIBUTING.md), as well as our [code of conduct](https://avanade.github.io/code-of-conduct/).

If you're happy to follow these guidelines, then check out the [getting started](./docs/start-here.md) guide.

## Who is Avanade?

[Avanade](https://www.avanade.com) is the leading provider of innovative digital and cloud services, business solutions and design-led experiences on the Microsoft ecosystem, and the power behind the Accenture Microsoft Business Group.