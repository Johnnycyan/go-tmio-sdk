<p align="center">
  <img src="https://raw.githubusercontent.com/Johnnycyan/Twitch-APIs/main/OneMoreDayIcon.svg" width="100" alt="project-logo">
</p>
<p align="center">
    <h1 align="center">GO-TMIO-SDK</h1>
</p>
<p align="center">
    <em>Accelerate your trackmania coding journey with ease.</em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/last-commit/Johnnycyan/go-tmio-sdk?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/Johnnycyan/go-tmio-sdk?style=default&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/Johnnycyan/go-tmio-sdk?style=default&color=0080ff" alt="repo-language-count">
	<img src="https://goreportcard.com/badge/github.com/Johnnycyan/go-tmio-sdk" alt="go-report-card">
<p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>

<hr>

##  Overview

The go-tmio-sdk project is a robust open-source toolkit designed to streamline interactions with the Trackmania API. It offers functionalities to retrieve and cache player information efficiently, including player IDs and names. With a focus on error handling and secure environment variable management, the project enhances the overall experience of working with the TMIO API. The go-tmio-sdk project stands out for its capability to simplify data fetching, caching, and retrieval tasks, catering to developers looking to integrate Trackmania API functionalities seamlessly.

---

##  Getting Started

**System Requirements:**

* **Internet**

###  Setting it Up

> 1. go get the package:
> ```console
> $ go get github.com/Johnnycyan/go-tmio-sdk
> ```
> 2. Import the package in your go file
> ```
> import (
>   tmio "github.com/Johnnycyan/go-tmio-sdk"
> )
> ```
> 3. Create a .env file for your project with:
> ```
> NAME=<your-username>
> ```

###  Usage

> Get a Player ID from a Player Name
> ```
> username := "username"
> playerID, err := tmio.GetPlayerID(username)
> if err != nil {
>   log.Println(err)
>   return
> }
> ```
> Get a Formatted Player Name from a non-formatted Player Name
> ```
> username := "username"
> playerFormatted, err := tmio.GetFormattedName(username)
> if err != nil {
>   log.Println(err)
>   return
> }
> ```

---
