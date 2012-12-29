(function () {
    $(function () {
        // Ideas from gnuplot's quantize_normal_tic
        //
        // returns between 5 and 10 tics
        function axisTics(data) {
            var min = data[0],
                max = data[data.length - 1],
                range = max - min,
                size,
                axis_min, axis_max,
                ret = [];
            // Ensure min and max don't crash into the axis.
            min -= range * 0.05;
            max += range * 0.05;
            size = Math.pow(10, Math.floor(Math.log(range) / Math.log(10)));

            // By construction 1 <= range / size < 10 as what we're doing is:
            //     range / 10 ** floor(log10(range))
            // Having a graph with fewer than 5 tics is not great, so we scale the axis if
            // this would happen. We use factors of 2 and 5 because humans are used to
            // them.
            if (range / size < 2) {
                // [0, 2, 4, 6]
                size *= 0.2;

            } else if (range / size < 5) {
                // [0, 5, 10, 15]
                size *= 0.5;

            } else {
                // [0, 1, 2, 3]
                size *= 1;
            }

            axis_min = min - (min % size);
            axis_max = max + size - (max % size);

            for (var i = axis_min; i <= axis_max; i += size) {
                ret.push(i);
            }

            return ret;
        }

        function logAxisTicks(data) {
            var min = data[0],
                max = data[data.length - 1],

                axis_min = Math.floor(Math.log(min) / Math.log(10)),
                axis_max = Math.ceil(Math.log(max) / Math.log(10)),
                pattern,
                ret = [];

            if (axis_max - axis_min < 1) {
                pattern = [1, 2, 3, 4, 5, 6, 7, 8, 9];

            } else if (axis_max - axis_min < 4) {
                pattern = [1, 2, 4, 6, 8];

            } else if (axis_max - axis_min < 8) {
                pattern = [1, 2.5, 5, 7.5];

            } else {
                pattern = [1, 5];

            }

            for (var i = axis_min; i < axis_max; i += 1) {
                ret = ret.concat(pattern.map(function (t) {
                    return Math.pow(10, i) * t;
                }));
            }

            ret.push(Math.pow(10, axis_max));

            return ret;
        }

        var height = 400, width = 600;

        window.raw = function (data) {

            data.sort(d3.ascending);
            var axis = axisTics(data);
            axis[0] = Math.max(axis[0], 1);

            var x = d3.scale.linear().domain([0, data.length - 1]).range([40, width]);
            var px = d3.scale.linear().domain([0, 100]).range([40, width]);
            var y = d3.scale.linear().domain([d3.min(axis), d3.max(axis)]).range([height - 20, 20]);

            var viz = d3.select('#chart').style('width', width).style('height', height)
                        .append('svg:svg').attr('width', width).attr('height', height).attr('class', 'viz')
                        .append('svg:g');

            viz.selectAll('path.line').data([data]).enter().append('svg:path')
                .attr('d', d3.svg.line().x(function (d, i) { return x(i); }).y(y));

            viz.selectAll('.tickx').data([0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100]).enter().append('svg:line')
                .attr('y1', y(axis[0])).attr('y2', y(axis[0]) - 3).attr('x1', px).attr('x2', px);

            viz.selectAll('.ticky').data(axis).enter().append('svg:line')
                .attr('x1', x(0)).attr('x2', x(0) + 3).attr('y1', y).attr('y2', y);

            viz.selectAll('.labely').data(axis).enter().append('svg:text')
                .text(function (d) { return d; }).attr('x', 0).attr('y', function (d) { return y(d) + 5; });

        };

        function mode(kde) {
            var max = -Infinity,
                ret;

            kde.forEach(function (d) {
                if (d[1] > max) {
                    max = d[1];
                    ret = d;
                }
            });

            return ret;
        }

        function expectation(kde) {
            var accum = 0,
                a, b;

            for (var i = 0; i < kde.length - 1; i++) {
                a = kde[i];
                b = kde[i + 1];

                accum += (b[0] - a[0]) * (a[0] * a[1] + b[0] * b[1]) / 2;
            }

            return accum;
        }

        function percentiles(data) {
            return [0, 0.01, 0.1, 0.25, 0.5, 0.75, 0.9, 0.99, 1].map(function (i) {
                return [i, data[Math.ceil((data.length - 1) * i)]];
            });
        }

        function mean(data) {
            var accum = 0;

            data.forEach(function (x) {
                accum += x;
            });

            return accum / data.length;
        }

        // TODO: figure out http://www.umiacs.umd.edu/labs/cvl/pirl/vikas/Software/optimal_bw/optimal_bw_code.htm
        function cdf(data) {
            var axisTicks = logAxisTicks(data);
            var x = d3.scale.log().domain([axisTicks[0], axisTicks[axisTicks.length - 1]]).range([40, width]);
            var y = d3.scale.linear().domain([1, 0]).range([40, height]);
            var kde = science.stats.kde().sample(data.map(function (x) { return Math.log(x); }))(d3.range(Math.log(data[0]), Math.log(data[data.length - 1]), 0.05));
            var yk = d3.scale.linear().domain([d3.max(kde.map(function (x) { return x[1]; })), d3.min(kde.map(function (x) { return x[1]; }))]).range([40, height]);

            var viz = d3.select('#cdf').style('width', width).style('height', height)
                        .append('svg:svg').attr('width', 40 + width).attr('height', height + 40).attr('class', 'viz')
                        .append('svg:g');

            var cumulative = [], i = 0;

            while (i < data.length) {
                while (data[i + 1] === data[i]) {
                    i += 1;
                }
                cumulative.push([data[i], (i + 1) / data.length]);
                i += 1;
            }

            viz.selectAll('line.ydivisions.strong').data([0, 1]).enter().append('svg:line')
                .attr('class', 'ydivisions strong')
                .attr('x1', x(axisTicks[0])).attr('x2', x(axisTicks[axisTicks.length - 1]))
                .attr('y1', y).attr('y2', y);

            viz.selectAll('line.ydivisions.weak').data(d3.range(0, 1, 0.25)).enter().append('svg:line')
                .attr('class', 'ydivisions weak')
                .attr('x1', x(axisTicks[0])).attr('x2', x(axisTicks[axisTicks.length - 1]))
                .attr('y1', y).attr('y2', y);

            viz.selectAll('line.xdivisions').data(axisTicks).enter().append('svg:line')
                .attr('class', function (d) {
                    if (d.toString().match(/^10*$/)) {
                        return 'xdivisions strong';
                    } else {
                        return 'xdivisions weak';
                    }
                })
                .attr('y1', y(0)).attr('y2', y(1))
                .attr('x1', x).attr('x2', x);

            viz.selectAll('text.xticklabels').data(axisTicks).enter().append('svg:text')
                .attr('class', 'xticklabels')
                .text(function (d) {
                    if (d.toString().match(/^10*$/)) {
                        return d;
                    } else {
                        return "";
                    }
                })
                .attr('x', x)
                .attr('y', height + 15);

            viz.selectAll('text.xtitle').data([0]).enter().append('svg:text')
                .attr('class', 'xtitle')
                .text("Email size (kb)")
                .attr('x', 40 + width / 2)
                .attr('y', height + 30);

            viz.selectAll('text.ytitle').data([0]).enter().append('svg:text')
                .attr('class', 'ytitle')
                .text('Probability')
                .attr('transform', 'rotate(270 30 ' + (40 + height / 2) + ')')
                .attr('x', 20)
                .attr('y', 40 + height / 2);

            viz.selectAll('text.y2title').data([0]).enter().append('svg:text')
                .attr('class', 'y2title')
                .text('Count so far')
                .attr('transform', 'rotate(90 ' + (width + 10) + ' ' + (40 + height / 2) + ')')
                .attr('x', width + 10)
                .attr('y', 40 + height / 2);

            // Ensure kde hits axis at both ends.
            kde.unshift([kde[0][0], 0]);
            kde.push([kde[kde.length - 1][0], 0]);

            viz.selectAll('path.pdf').data([[[kde[0][0], 0]].concat(kde)]).enter().append('svg:path')
                .attr('class', 'pdf')
                .attr('d', d3.svg.line().x(function (d) {
                    return x(Math.exp(d[0]));
                }).y(function (d) {
                    return yk(d[1]);
                }));

            viz.selectAll('path.cdf').data([cumulative]).enter().append('svg:path')
                .attr('class', 'cdf')
                .attr('d', d3.svg.line().x(function (d) {
                    return x(d[0]);
                }).y(function (d) { return y(d[1]); }));


            var m = mode(kde);

            viz.selectAll('circle.mode').data([mode(kde)]).enter().append('svg:circle')
                .attr('class', 'mode')
                .attr('r', 4)
                .attr('cx', x(Math.exp(m[0])))
                .attr('cy', yk(m[1]));

            viz.selectAll('circle.expectation').data([expectation(kde)]).enter().append('svg:circle')
                .attr('class', 'expectation')
                .attr('r', 4)
                .attr('cx', function (d) { return x(Math.exp(d)); })
                .attr('cy', y(0.25));

//            viz.selectAll('circle.foo').data([mean(data)]).enter().append('svg:circle')
//                .attr('r', 4)
//                .attr('cx', x(mean(data)))
//                .attr('cy', y(0.25));

            viz.selectAll('circle.percentile').data(percentiles(data)).enter().append('svg:circle')
                .attr('class', 'percentile')
                .attr('r', 4)
                .attr('cx', function (d) { return x(d[1]); })
                .attr('cy', function (d) { return y(d[0]); });
        }

        window.emails = window.emails.filter(function (x) { return x !== 0; });

        window.normal = window.normal.filter(function (x) { return x !== 0; }).sort(d3.ascending);

        cdf(window.emails.map(Number));
//        cdf(window.files.map(function (x) { return x / 1024; }));
    });
}).call(this);
