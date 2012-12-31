(function () {
    $(function () {
        var height = 400, width = 600;

        function cdf(data) {
            var kde = science.stats.distribution.kde().sample(data).resolution(200);

            var x;
            if (kde.feelsLogarithmic()) {
                kde.log(true);
                x = d3.scale.log();
            } else {
                x = d3.scale.linear();
            }
            x.clamp(true).domain([data[0], data[data.length - 1]]).range([40, width]).nice();
            var axisTicks = x.ticks(10);
            var y = d3.scale.linear().domain([1, 0]).range([40, height]);
            var yk = d3.scale.linear().domain([kde.max(), 0]).range([40, height]);

            var viz = d3.select('#cdf').style('width', width).style('height', height)
                        .append('svg:svg').attr('width', 40 + width).attr('height', height + 40).attr('class', 'viz')
                        .append('svg:g');

            viz.selectAll('line.ydivisions.strong').data([0, 1]).enter().append('svg:line')
                .attr('class', 'ydivisions strong')
                .attr('x1', x(axisTicks[0])).attr('x2', x(axisTicks[axisTicks.length - 1]))
                .attr('y1', y).attr('y2', y);

            viz.selectAll('line.ydivisions.weak').data(d3.range(0, 1, 0.25)).enter().append('svg:line')
                .attr('class', 'ydivisions weak')
                .attr('x1', x(axisTicks[0])).attr('x2', x(axisTicks[axisTicks.length - 1]))
                .attr('y1', y).attr('y2', y);

            viz.selectAll('line.xdivisions').data(axisTicks).enter().append('svg:line')
                .attr('class', function (d, i) {
                    if (i === 0 || d.toString().match(/^10*$/) || i === axisTicks.length - 1) {
                        return 'xdivisions strong';
                    } else {
                        return 'xdivisions weak';
                    }
                })
                .attr('y1', y(0)).attr('y2', y(1))
                .attr('x1', x).attr('x2', x);

            viz.selectAll('text.xticklabels').data(axisTicks).enter().append('svg:text')
                .attr('class', 'xticklabels')
                .text(function (d, i) {
                    if (i === 0 || d.toString().match(/^10*$/) || i === axisTicks.length - 1) {
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

            viz.selectAll('path.pdf').data([kde.pdf()]).enter().append('svg:path')
                .attr('class', 'pdf')
                .attr('d', d3.svg.line().x(function (d) {
                    return x(d[0]);
                }).y(function (d) {
                    return yk(d[1]);
                }));

            viz.selectAll('path.cdf').data([kde.cdf()]).enter().append('svg:path')
                .attr('class', 'cdf')
                .attr('d', d3.svg.line().x(function (d) {
                    return x(d[0]);
                }).y(function (d) { return y(d[1]); }));

            viz.selectAll('circle.mode').data([kde.mode()]).enter().append('svg:circle')
                .attr('class', 'mode')
                .attr('r', 4)
                .attr('cx', x)
                .attr('cy', function (d) { return yk(kde(d)); });

            viz.selectAll('circle.expectation').data([kde.expectation()]).enter().append('svg:circle')
                .attr('class', 'expectation')
                .attr('r', 4)
                .attr('cx', x)
                .attr('cy', function () { return y(0.25); });

            viz.selectAll('circle.percentile').data(kde.qf()).enter().append('svg:circle')
                .attr('class', 'percentile')
                .attr('r', 4)
                .attr('cx', function (d) { return x(d[0]); })
                .attr('cy', function (d) { return y(d[1]); });

            viz.selectAll('line.fugi').data([0]).enter().append('svg:line')
                .attr('class', 'fugi')
                .attr('y1', y(0)).attr('y2', y(1))
                .attr('x1', x).attr('x2', x);

            viz.selectAll('circle.fugired').data([0]).enter().append('svg:circle')
                .attr('class', 'fugired')
                .attr('r', 5)
                .attr('cx', 0)
                .attr('cy', 0);

            viz.selectAll('circle.fugiblue').data([0]).enter().append('svg:circle')
                .attr('class', 'fugiblue')
                .attr('r', 5)
                .attr('cx', 0)
                .attr('cy', 0);

            var $svg = $(viz[0][0]).closest('svg');

            $svg.mousemove(function (e) {
                var pos = x.invert(e.clientX - 10);
                if (pos < axisTicks[0]) {
                    pos = axisTicks[0];
                } else if (pos > axisTicks[axisTicks.length - 1]) {
                    pos = axisTicks[axisTicks.length - 1];
                }
                $svg.find('line.fugi').attr('x1', x(pos)).attr('x2', x(pos));
                $svg.find('.fugired, .fugiblue').attr('cx', x(pos));
                $svg.find('.fugiblue').attr('cy', yk(kde(pos)));
                $svg.find('.fugired').attr('cy', y(kde.inverseQuantile(pos)));
            }).mouseover(function () {
                $svg.find('line.fugi').show();
            }).mouseout(function () {
                $svg.find('line.fugi').hide();
            });
        }

        window.emails = window.emails.filter(function (x) { return x !== 0; });

        window.normal = window.normal.map(Number).filter(function (x) { return x !== 0; }).sort(d3.ascending);

        window.faithful = window.faithful.sort(d3.ascending);
//        cdf(window.faithful);
        window.foo = [1, 2];
        cdf(window.emails);
        cdf(window.normal);
        cdf(window.files);
        cdf(window.faithful);
        cdf(window.foo);
    });
}).call(this);
