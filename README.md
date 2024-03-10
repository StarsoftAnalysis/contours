# Contours 
![example workflow](https://github.com/starsoftanalysis/contours/actions/workflows/WORKFLOW-FILE/badge.svg)
This is a basic script for pulling out contours/regions from an image.

It uses Moore Neighbourhood contour tracing. For more details on the algorithm/approach see [this site](http://www.imageprocessingplace.com/downloads_V3/root_downloads/tutorials/contour_tracing_Abeer_George_Ghuneim/moore.html).

```js
const output = contours(imageData)

// list of matching contours
```

* imageData should be thresholded (any non-zero pixel is counted as _on_)
* only the red value is the only one used

## Development

```bash
yarn install #or, npm

# build dist/contours
npm run build

# run tests
npm test

# run tests on file change
npm run test:watch
```

## Status

Forked from [Ben Foxall's project](https://github.com/benfoxall/contours) in March 2024.
