// TODO: figure out http://www.umiacs.umd.edu/labs/cvl/pirl/vikas/Software/optimal_bw/optimal_bw_code.htm
science.stats.distribution.kde = function () {

    var underlying = science.stats.kde(),
        sample = [],
        cache = {},
        log = false,
        resolution = 100;

    function delog(x) {
        return log ? Math.exp(x) : x;
    }

    function enlog(x) {
        return log ? Math.log(x) : x;
    }

    function logSample() {
        if (!cache.log) {
            cache.log = log ? sample.map(function (x) {
                return Math.log(x);
            }) : sample;
        }
        return cache.log;
    }

    // Return the KDE as an array of [x, probability-density] pairs.
    function quantized() {
        if (!cache.quantized) {
            var first = logSample()[0],
                last  = logSample()[logSample().length - 1],
                calculated = underlying(d3.range(first, last, (last - first) / resolution));

            // Ensure KDE function hits the axis.
            calculated.shift([first, 0]);
            calculated.push([last, 0]);

            cache.quantized = calculated;
        }
        return cache.quantized;
    }

    function delogKde() {
        if (log) {
            return quantized().map(function (d) { return [Math.exp(d[0]), d[1]]; });
        } else {
            return quantized();
        }
    }

    function kde(x) {
        return underlying([enlog(x)])[0][1];
    }

    kde.bandwidth = function (x) {
        if (!arguments.length) { return underlying.bandwidth(); }
        underlying.bandwidth(x);
        return kde;
    };

    kde.kernel = function (x) {
        if (!arguments.length) { return underlying.kernel(); }
        underlying.kernel(x);
        return kde;
    };

    /* The samples from which to derive the KDE.
     *
     * An array of numbers.
     */
    kde.sample = function (x) {
        if (!arguments.length) { return sample; }

        sample = x;
        cache = {};
        underlying.sample(logSample());
        return kde;
    };

    /* Whether to perform KDE on the log of the samples instead of the samples themselves
     *
     * The main reason to set this is that the symmetrical kernel used by default will
     * give very bizarre results if the data distribution is not locally symmetrical,
     * however it also interacts with the resolution parameter so that the resolution will
     * appear uniform when plotted on a log-graph.
     *
     * If you set this you should plot the KDE function on a log-graph (and vice-versa).
     *
     * [1] http://www.ebyte.it/library/docs/math04a/PdfChangeOfCoordinates04.html
     */
    kde.log = function (x) {
        if (!arguments.length) { return log; }

        log = x;
        cache = {};
        underlying.sample(logSample());
        return kde;
    };

    /* How many samples will be used.
     *
     * This effects the number of points returned by kde() with no arguments,
     * and (to a small amount) the accuracy of the estimate provided by .expectation().
     *
     * You probably want to set this to the same order of magnitude as the number of
     * pixels in the graph so that the curve looks smooth.
     */
    kde.resolution = function (x) {
        if (!arguments.length) { return resolution; }

        resolution = x;
        cache = {};
        return kde;
    };

    /* The 'mode' of the kde.
     *
     * This is the highest point of the probability density function; intuitively the value
     * around which samples tend to cluster. (e.g. most files are 4kb)
     *
     * (The more usual definition of mode as 'the value that occurs most often' falls
     * apart a bit when you have continuous data as the chances of the same exact value
     * occuring more than once is negligable).
     *
     * http://en.wikipedia.org/wiki/Mode_(statistics)
     */
    kde.mode = function () {
        var max = -Infinity,
            kde = quantized(),
            ret;

        kde.forEach(function (d) {
            if (d[1] > max) {
                max = d[1];
                ret = d[0];
            }
        });

        return delog(ret);
    };

    /* The 'expectation' of the kde.
     *
     * This is the first moment of the data; and corresponds to most people's intuitive
     * value for 'average'. (e.g. a typical file is 32kb, size = 32kb * number of files)
     *
     * (The arithmetic mean is often used in place of expectation, as they're roughly
     * equivalent for symmetrical distributions. The expectation degrades more gracefully
     * when the data is very skewed).
     *
     * http://en.wikipedia.org/wiki/Expected_value
     */
    kde.expectation = function () {
        var accum = 0,
            kde = quantized(),
            a, b;

        for (var i = 0; i < kde.length - 1; i++) {
            a = kde[i];
            b = kde[i + 1];

            // integrate(x * f(x))
            // => sum the area of all quadrilaterals of size d
            accum += (b[0] - a[0]) * (a[0] * a[1] + b[0] * b[1]) / 2;
        }

        return delog(accum);
    };

    /* The 'median' of the data.
     *
     * Half the data is smaller than the median, the other half is bigger.
     * (e.g. half of all files are smaller than 16kb)
     */
    kde.median = function () {
        return kde.percentile(0.5);
    };

    kde.max = function () {
        return kde(kde.mode());
    };

    kde.quantile = function (x) {
        return science.stats.quantiles(sample, [x])[0];
    };

    kde.inverse_quantile = function (x) {
        return science.bisect(sample, function (d) {
            return d < x;
        }) / sample.length;
    };

    // The probability density function as an array of [x, y] pairs.
    kde.pdf = function () {
        return delogKde();
    };

    // Some quartiles as an array of [x, y] pairs.
    kde.qf = function () {
        var ys = [0, 0.01, 0.1, 0.25, 0.5, 0.75, 0.9, 0.99, 1];

        return ys.map(function (y) {
            return [kde.quantile(y), y];
        });
    };

    // The cumulative probability function as an array of [x, y] pairs
    kde.cdf = function () {
        if (!cache.cdf) {
            var cumulative = [], i = 0;
            while (i < sample.length) {
                while (sample[i + 1] === sample[i]) {
                    i += 1;
                }
                cumulative.push([sample[i], (i + 1) / sample.length]);
                i += 1;
            }
            cache.cdf = cumulative;
        }
        return cache.cdf;
    };

    return kde;

};
