# Moore-Neighbourhood Contours 
![example workflow](https://github.com/starsoftanalysis/contours/actions/workflows/WORKFLOW-FILE/badge.svg)

This is a Go programme for pulling out contours/regions from an image.

It uses Moore-Neighbourhood contour tracing. For more details on the algorithm/approach see [this site](http://www.imageprocessingplace.com/downloads_V3/root_downloads/tutorials/contour_tracing_Abeer_George_Ghuneim/moore.html).

## Status

Forked from [Ben Foxall's project](https://github.com/benfoxall/contours) in March 2024,
and translated from Javascript to Go.

### Issues

* A non-closed thin shape such as test7.png has to look at the shape from both sides to find the whole contour,
  but then ends up with repeated pixels or, if those are excluded, pixels out of order.
