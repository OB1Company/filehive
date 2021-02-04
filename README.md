[![Contributors][contributors-shield]][contributors-url] 
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/OB1Company/filehive">
    <img src="https://filehive.app/filehive-logo.png" alt="Logo">
  </a>

  <h3 align="center">Filehive</h3>

  <p align="center">
    A Filecoin-backed marketplace for datasets!
    <br />
    <a href="https://github.com/OB1Company/filehive"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://beta.filehive.app">View Beta</a>
    ·
    <a href="https://github.com/OB1Company/filehive/issues">Report Bug</a>
    ·
    <a href="https://github.com/OB1Company/filehive/issues">Request Feature</a>
  </p>
</p>



<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

Data is the lifeblood of the Internet and therefore one of the most valuable digital assets that can be exchanged. Filehive is a new kind of marketplace built on top of Filecoin and IPFS to provide everyone across the planet an easy and cheap way to distribute and consume diverse datasets. 

The OB1 team has been working with Protocol Labs for many years now and this is the first Filecoin first project to come out of their development shop.

### Built With

* [IPFS](https://ipfs.io)
* [Filecoin](https://filecoin.io)
* [Powergate](https://github.com/textileio/powergate)
* [React](https://reactjs.org)
* [Golang](https://golang.org)

<!-- GETTING STARTED -->
## Getting Started

Follow the instructions below to get your Filehive environment working.

### Prerequisites

Filehive requires three components:
* [Textile's Powergate](#textile-powergate)
* [Filehive Go API Server](#filehive-go-api-server)
* [Filehive React Web Application](#filehive-ui)

#### Textile Powergate

See Textile's amazing instructions for setting a Powergate server. If you would like to run this application locally in development mode you can use a `localnet` version of Powergate. Textile provides details on how to use Docker to set it all up very quickly. https://docs.textile.io/powergate/localnet/

#### Filehive Go API Server

Once you have your Powergate server up and running or you have access to a hosted instance that you can connect to from Filehive you can proceed to installing and starting the Filehive Go API Server.

Make sure you have Go 1.15 or above to run this server.

1. Go get this project on your machine (`go get -u https://github.com/OB1Company/filehive`)
2. Change into your source code folder (`cd $GOPATH/src/github.com/OB1Company/filehive`)
3. Start the server (`go run main.go`)

Once you start the server a Filehive data repository will be created on your machine in your OS-specific location. If you need to customize your configuration to connect to a different Powergate server there is a `filehive.conf` file in your data repository folder that you can modify. Once you've updated your conf file you need to restart the server for changes to take effect.

#### Filehive UI

Filehive is a mono repo so the Go repo has the React web UI included (in the /web folder). You can either clone this repo in a different location or use the existing Go repo and navigate to the `/web` folder.

1. Change into the `/web` folder
2. Install dependencies using `yarn` or `npm` (`yarn install`)
3. Start the node application (`yarn start`)

You can customize your application by specifying settings in the `.env` file in the root of the `/web` folder. You can see a [sample `.env` file](https://github.com/OB1Company/filehive/blob/master/web/.env_sample).

1. Create `.env` file in root `/web` folder
2. Update contents of `.env`
3. Restart the node application (`yarn start`)

The application will start by default at [http://localhost:3000](http://localhost:3000)

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/OB1Company/filehive/issues) for a list of proposed features (and known issues).



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Filehive - [@TheFilehive](https://twitter.com/thefilehive)

OB1 Company: [https://github.com/OB1Company](https://github.com/OB1Company)



<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements
* [Protocol Labs](https://protocol.ai)
* [Textile](https://textile.io)


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/OB1Company/filehive.svg?style=for-the-badge
[contributors-url]: https://github.com/OB1Company/filehive/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/OB1Company/filehive.svg?style=for-the-badge
[forks-url]: https://github.com/OB1Company/filehive/network/members
[stars-shield]: https://img.shields.io/github/stars/OB1Company/filehive.svg?style=for-the-badge
[stars-url]: https://github.com/OB1Company/OB1Company/filehive
[issues-shield]: https://img.shields.io/github/issues/OB1Company/filehive.svg?style=for-the-badge
[issues-url]: https://github.com/OB1Company/filehive/issues
[license-shield]: https://img.shields.io/github/license/OB1Company/filehive.svg?style=for-the-badge
[license-url]: https://github.com/OB1Company/filehive/blob/master/LICENSE.txt
[product-screenshot]: images/screenshot.png
