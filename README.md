<div align="center" style="text-align: center;">
  <img src="gopher-batman.png" alt="Description" style="width: 50%; max-width: 300px;">
</div>

<p align="center"><a href="https://pkg.go.dev/github.com/Permify/policy-enforcer?tab=doc"
target="_blank"></a><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="go version" />&nbsp;&nbsp;<img src="https://img.shields.io/github/license/grindlemire/gothem-stack?style=for-the-badge" alt="license" />

# gothem-stack

An end to end [htmx](https://htmx.org) and [go-templ](https://templ.guide) template using [echo](https://echo.labstack.com/) for the web server and [mage](https://magefile.org/) for deployment. Other than that this is unopinionated so bring all your other favorite technologies :)

**Go**\
**T**empl\
**H**TMX\
**E**cho\
**M**age

For the frontend libraries it uses [tailwindcss](https://tailwindcss.com/) and [daisyUI](https://daisyui.com/) for styling. You can easily integrate [alpine-js](https://alpinejs.dev/) or [hyperscript](https://hyperscript.org/) if you desire. The example integrates with both.

This template uses my [graft](https://github.com/grindlemire/graft) dependency injection library to bootstrap up the server. If you would prefer a non bootstrapped version checkout the `v1` branch or tag.

## Quickstart

1. Install npm in your path `brew install node; brew install npm`
1. Install mage in your path `brew install mage`. See [https://magefile.org/](https://magefile.org/) for other installation instructions
1. Run `mage install`
1. Run `mage run`
1. Open <http://localhost:7331> (the server is listening on :4433 but templ injects a watcher with an autoreload script)

When you make changes to your templ files or any of your go code everything will regenerate and then autoreload your web page

## Basic Commands

`mage run` - Run an interactive development environment that will automatically reload on any file change. Listens on port :4433 and has an autoreload page on :7331

`mage install` - Install all the dependencies

`mage templ` - Do a one time regeneration of your templ files

`mage build` - Do a one time build of the go files

`mage tidy` - Run `go mod tidy`

## Versions

- **v2.x** (current `main`) Uses my [graft](https://github.com/grindlemire/graft) library for managing the node dependencies
- **v1.x** (`v1` branch / `v1.0.0` tag) Is a less opinionated version without a DI framework

To use v1:
git clone --branch v1 <https://github.com/grindlemire/gothem-stack.git>

## Cloud Deployment (Optional)

The project also includes optional support for deploying your code to Google Cloud Run and Firebase Hosting. **This is by no means required to use this project, if you choose to you can just ignore all these commands and use it without cloud integration**. To use these features, you'll need to set up Google Cloud and Firebase projects first.

See the example running at [gothem-stack.web.app](https://gothem-stack.web.app/)

### Cloud Deployment Commands

`mage auth [gcloud|firebase]` - Authenticate with Google Cloud or Firebase services

`mage deploy [backend|frontend|all]` - Deploy your application

- `backend` - Deploys to Google Cloud Run
- `frontend` - Deploys to Firebase Hosting
- `all` - Deploys both (default if no argument provided)

`mage release [backend|frontend|all]` - Move to a specific version of the backend service

`mage destroy` - Tear down all cloud infrastructure created during the deployment

### Required Environment Variables for Cloud Deployment

If using cloud deployment, you'll need to set these environment variables:

- `GOOGLE_CLOUD_PROJECT` - Your Google Cloud project ID
- `GOOGLE_CLOUD_REGION` - The region to deploy to (defaults to us-central1)
- `CLOUD_RUN_SERVICE_NAME` - The name of your Cloud Run service

You can place these variables in a env file to declaratively use them. For example
you could run

```bash
echo "GOOGLE_CLOUD_PROJECT=foobar-123\nGOOGLE_CLOUD_REGION=us-central1\nCLOUD_RUN_SERVICE_NAME=gothem-stack" > gothem-stack.env
```

Then you can run `ENV=gothem-stack mage deploy` to use the configured environment variables.

### But Gotham is spelled with an 'a'

Yea I know it's spelled with an 'a' :] I was trying to come up with a name that was easier to say and interact with than 'htmx-templ-template' and gothem is what I came up with.
