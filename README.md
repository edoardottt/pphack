<h1 align="center">
  pphack
  <br>
</h1>

<h4 align="center">The Most Advanced Prototype Pollution Scanner</h4>

<h6 align="center"> Coded with 💙 by edoardottt </h6>

<p align="center">

  <a href="https://github.com/edoardottt/pphack/actions">
      <img src="https://github.com/edoardottt/pphack/actions/workflows/go.yml/badge.svg" alt="go action">
  </a>

  <a href="https://goreportcard.com/report/github.com/edoardottt/pphack">
      <img src="https://goreportcard.com/badge/github.com/edoardottt/pphack" alt="go report card">
  </a>

<br>
  <!--Tweet button-->
  <a href="https://twitter.com/intent/tweet?text=pphack%20-%20The%20Most%20Advanced%20Prototype%20Pollution%20Scanner%20https%3A%2F%2Fgithub.com%2Fedoardottt%2Fpphack%20%23golang%20%23github%20%23linux%20%23infosec%20%23bugbounty" target="_blank">Share on Twitter!
  </a>
</p>

<p align="center">
  <a href="#install-">Install</a> •
  <a href="#get-started-">Get Started</a> •
  <a href="#examples-">Examples</a> •
  <a href="#changelog-">Changelog</a> •
  <a href="#contributing-">Contributing</a> •
  <a href="#license-">License</a>
</p>

<p align="center">
  <img src="https://github.com/edoardottt/images/blob/main/pphack/pphack.gif">
</p>
  
Install 📡
----------

### Using Snap

```console
sudo snap install pphack
```

### Using Go

```console
go install github.com/edoardottt/pphack/cmd/pphack@latest
```

Get Started 🎉
----------

```console
Usage:
  pphack [flags]

Flags:
INPUT:
   -u, -url string   Input URL
   -l, -list string  File containing input URLs

CONFIGURATIONS:
   -c, -concurrency int  Concurrency level (default 50)
   -t, -timeout int      Connection timeout in seconds (default 10)

OUTPUT:
   -o, -output string  File to write output results
   -v, -verbose        Verbose output
   -s, -silent         Silent output. Print only results
```

Examples 💡
----------

Scan a single URL

```console
pphack -u https://edoardottt.github.io/pp-test/
```

```console
echo https://edoardottt.github.io/pp-test/ | pphack
```

Scan a list of URLs

```console
pphack -l targets.txt
```

```console
cat targets.txt | pphack
```

Changelog 📌
-------

Detailed changes for each release are documented in the [release notes](https://github.com/edoardottt/pphack/releases).

Contributing 🛠
-------

Just open an [issue](https://github.com/edoardottt/pphack/issues) / [pull request](https://github.com/edoardottt/pphack/pulls).

Before opening a pull request, download [golangci-lint](https://golangci-lint.run/usage/install/) and run

```bash
golangci-lint run
```

If there aren't errors, go ahead :)

License 📝
-------

This repository is under [MIT License](https://github.com/edoardottt/pphack/blob/main/LICENSE).  
[edoardoottavianelli.it](https://www.edoardoottavianelli.it) to contact me.
