# TODO: figure out http://www.umiacs.umd.edu/labs/cvl/pirl/vikas/Software/optimal_bw/optimal_bw_code.htm
science.stats.distribution.kde = ->

  underlying = science.stats.kde()
  sample = []
  cache = {}
  scale = d3.scale.identity()
  resolution = 100

  scaledSample = ->
    sample.map scale

  # Return the KDE as an array of [x, probability-density] pairs.
  scaledPdf = ->
    return cache.scaledPdf if cache.scaledPdf

    scaled = scaledSample()
    first  = scaled[0]
    last   = scaled[scaled.length - 1]
    step   = (last - first) / resolution

    calculated = underlying(d3.range(first, last, step))

    # Ensure KDE function hits the axis.
    # (this is necessary to calculate expectation correctly)
    while calculated[0][1] > 1e-3
      first -= step
      calculated.unshift underlying([first])[0]

    while calculated[calculated.length - 1][1] > 1e-3
      last += step
      calculated.push underlying([last])[0]

    calculated.unshift [first, 0]
    calculated.push [last, 0]

    cache.scaledPdf = calculated

  kde = (x) ->
    underlying([scale(x)])[0][1]

  kde.bandwidth = (x) ->
    return underlying.bandwidth() unless arguments.length
    underlying.bandwidth x
    kde

  kde.kernel = (x) ->
    return underlying.kernel() unless arguments.length
    underlying.kernel x
    kde


  # The samples from which to derive the KDE.
  kde.sample = (x) ->
    return sample  unless arguments.length
    sample = x
    cache = {}
    underlying.sample scaledSample()
    kde


  # The scale on which to perform kernel density estimation.
  #
  # This should be proportional to the scale with which you will display the
  # reuslting probability density function (as kernel density estimation is
  # essentially a formalization of the appearance of the data).
  #
  # In particular if you perform KDE on a linear axis and then display it on
  # a log graph, the symmetric gaussian kernel bleeds heavily to the left on
  # the very asymmetric axis.
  #
  # See http://www.ebyte.it/library/docs/math04a/PdfChangeOfCoordinates04.html
  kde.scale = (x) ->
    return scale  unless arguments.length
    scale = x
    cache = {}
    underlying.sample scaledSample()
    kde


  # How many samples will be used.
  #
  # This effects the number of points returned by kde() with no arguments,
  # and (to a small amount) the accuracy of the estimate provided by .expectation().
  #
  # You probably want to set this to the same order of magnitude as the number of
  # pixels in the graph so that the curve looks smooth.
  kde.resolution = (x) ->
    return resolution  unless arguments.length
    resolution = x
    cache = {}
    kde

  # The 'mode' of the kde.
  #
  # This is the highest point of the probability density function; intuitively the value
  # around which samples tend to cluster. (e.g. most files are 4kb)
  #
  # (The more usual definition of mode as 'the value that occurs most often' falls
  # apart a bit when you have continuous data as the chances of the same exact value
  # occuring more than once is negligable).
  #
  # http://en.wikipedia.org/wiki/Mode_(statistics)
  kde.mode = ->
    max = -Infinity
    ret = undefined

    scaledPdf().forEach (d) ->
      if d[1] > max
        max = d[1]
        ret = d[0]

    scale.invert ret

  # The 'expectation' of the kde.
  #
  # This is the first moment of the data; and corresponds to most people's intuitive
  # value for 'average'. (e.g. a typical file is 32kb, size = 32kb * number of files)
  #
  # (The arithmetic mean is often used in place of expectation, as they're roughly
  # equivalent for symmetrical distributions. The expectation degrades more gracefully
  # when the data is very skewed).
  #
  # http://en.wikipedia.org/wiki/Expected_value
  kde.expectation = ->
    accum = 0
    pdf = scaledPdf()
    i = 0

    while i < pdf.length - 1
      a = pdf[i]
      b = pdf[i + 1]

      # integrate(x * f(x))
      # => sum the area of all quadrilaterals of size d
      accum += (b[0] - a[0]) * (a[0] * a[1] + b[0] * b[1]) / 2
      i++

    scale.invert accum

  kde.mean = ->
    science.stats.mean(sample)

  kde.variance = ->
    science.stats.variance(sample)

  kde.size = ->
    sample.length

  # The 'median' of the data.
  #
  # Half the data is smaller than the median, the other half is bigger.
  # (e.g. half of all files are smaller than 16kb)
  kde.median = ->
    kde.quantile 0.5

  kde.max = ->
    kde kde.mode()

  kde.quantile = (q) ->
    index = q * (sample.length - 1)
    lo = Math.floor(index)
    hi = Math.ceil(index)
    if lo is hi or hi >= sample.length
      sample[lo]
    else
      sample[lo] + (index - lo) * (sample[hi] - sample[lo]) / (hi - lo)

  kde.inverseQuantile = (x) ->
    # sample[lo] < x <= sample[hi]
    hi = science.bisect(sample, (d) ->
      d < x
    )
    lo = hi - 1
    if hi is sample.length
      1
    else if hi is 0
      0
    else
      (lo + (hi - lo) * (x - sample[lo]) / (sample[hi] - sample[lo])) / (sample.length - 1)

  # The probability density function as an array of [x, y] pairs.
  kde.pdf = ->
      scaledPdf().map (d) ->
        [scale.invert(d[0]), d[1]]

  # Some quartiles as an array of [x, y] pairs.
  kde.qf = ->
    ys = [0, 0.01, 0.1, 0.25, 0.5, 0.75, 0.9, 0.99, 1]
    ys.map (y) ->
      [kde.quantile(y), y]

  # The cumulative probability function as an array of [x, y] pairs
  kde.cdf = ->
    unless cache.cdf
      cache.cdf = sample.map((x, i) ->
        [x, i / (sample.length - 1)]
      )
    cache.cdf

  kde.feelsLogarithmic = ->
    tenth = kde.quantile(0.01)
    fiftieth = kde.quantile(0.5)
    ninetieth = kde.quantile(0.99)

    # if 50% of the data takes up <10% of the graph, it's logarithmic
    if (fiftieth - tenth) / (ninetieth - tenth) < 0.2
      true
    else
      false

  kde
