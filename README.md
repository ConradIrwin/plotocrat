# Plotocrat - rule by graphs

Welcome to the code behind https://plotocrat.com, the simple
Kernel-Density-Estimation website that anyone can post datasets to.

### Running locally

- Get the code with `git clone git@github.com:ConradIrwin/plotocrat`
- Install Postgres (exercise for the reader, but `brew install postgres` on macOS)
- Create a user 'plotocrat' with a password 'plotocrat' and a database called 'plotocrat' that that user has full permissions to
- Run `./run.sh` which will use https://github.com/superhuman/lrt to live-reload on change.

### Deploying

- Install heroku and then `heroku login`
- Add the remote `heroku git:remote --app plotocrat2`
- Then push `git push heroku main`

## 2022-02-27 Notes

The current code is a bit of a mess, as the original version stopped working (it
had been last updated to the latest rails in 2016).

Rather than trying to upgrade to a newer rails, I rewrote the backend in Go, and
was luckily able to download the compiled javascript and CSS from
https://archive.org.

It would be nice to set up something like `esbuild` and modernize the javascript
from the `last-rails-release` tag. This will require:

- Converting coffeescript (v1) code to javascript
- Replacing the rails bundler with `esbuild` or similar
- Updating jQuery (or removing it) and upgrading d3.

## Meta-fu

Licensed under the MIT license. Bug-reports, and pull requests are welcome.
