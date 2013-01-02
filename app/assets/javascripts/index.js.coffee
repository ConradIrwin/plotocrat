(->
  $ ->
    cdf = (data) ->
      kde = science.stats.distribution.kde().sample(data).resolution(200)
      x = undefined
      if kde.feelsLogarithmic()
        d2 = []
        # TODO: figure out the correct solution...
        data.forEach (x) -> d2.push(x) unless x is 0
        data = d2
        kde = science.stats.distribution.kde().sample(data).resolution(200)
        x = d3.scale.log()
      else
        x = d3.scale.linear()

      kde.scale x.copy()
      x.clamp(true).domain([data[0], data[data.length - 1]]).range([40, width]).nice()
      axisTicks = x.ticks(5)

      y = d3.scale.linear().domain([1, 0]).range([40, height])
      yk = d3.scale.linear().domain([kde.max(), 0]).range([40, height])

      viz = d3.select("#cdf")
        .style("width", width)
        .style("height", height)
        .append("svg:svg")
        .attr("width", 40 + width)
        .attr("height", height + 40)
        .attr("class", "viz")
        .append("svg:g")

      viz.selectAll("line.ydivisions.strong").data([0, 1]).enter()
        .append("svg:line")
        .attr
          class: "ydivisions strong"
          x1: x(axisTicks[0])
          x2: x(axisTicks[axisTicks.length - 1])
          y1: y
          y2: y

      viz.selectAll("line.ydivisions.weak").data(d3.range(0, 1, 0.25)).enter()
        .append("svg:line")
        .attr
          class: "ydivisions weak"
          x1: x(axisTicks[0])
          x2: x(axisTicks[axisTicks.length - 1])
          y1: y
          y2: y

      viz.selectAll("line.xdivisions").data(axisTicks).enter()
        .append("svg:line")
        .attr
          class: (d, i) ->
            if i is 0 or d.toString().match(/^10*$/) or i is axisTicks.length - 1
              "xdivisions strong"
            else
              "xdivisions weak"
          y1: y(0)
          y2: y(1)
          x1: x
          x2: x

      viz.selectAll("text.xticklabels").data(axisTicks).enter()
        .append("svg:text").attr("class", "xticklabels")
        .text((d, i) ->
          raw = x.tickFormat(5)(d)

          raw.replace(/^1e\+([0-6])$/, (m, i) ->
            Math.pow(10, i)
          )

        ).attr
          x: x
          y: height + 15

      viz.selectAll("text.xtitle").data([0]).enter()
        .append("svg:text").attr("class", "xtitle")
        .text(d3.select('#cdf').attr('data-title'))
        .attr
          x: 40 + width / 2
          y: height + 30

      viz.selectAll("text.ytitle").data([0]).enter()
        .append("svg:text").attr("class", "ytitle")
        .text("Probability")
        .attr
          transform: "rotate(270 30 " + (40 + height / 2) + ")"
          x: 20
          y: 40 + height / 2

      viz.selectAll("text.y2title").data([0]).enter()
        .append("svg:text")
        .attr("class", "y2title")
        .text("Count so far")
        .attr
          transform: "rotate(90 " + (width + 10) + " " + (40 + height / 2) + ")"
          x: width + 10
          y: 40 + height / 2

      viz.selectAll("path.pdf").data([kde.pdf()]).enter()
        .append("svg:path")
        .attr
          class: "pdf"
          d: d3.svg.line().x((d) -> x d[0]).y((d) -> yk d[1])

      viz.selectAll("path.cdf").data([kde.cdf()]).enter()
        .append("svg:path")
        .attr
          class: "cdf"
          d: d3.svg.line().x((d) -> x d[0]).y((d) -> y d[1])

      viz.selectAll("circle.mode").data([kde.mode()]).enter()
        .append("svg:circle")
        .attr
          class: "mode"
          r: 4
          cx: x
          cy: (d) -> yk kde(d)

      viz.selectAll("circle.expectation").data([kde.expectation()]).enter()
        .append("svg:circle")
        .attr
          class: "expectation"
          r: 4
          cx: x
          cy: -> y 0.25

      viz.selectAll("circle.percentile").data(kde.qf()).enter()
        .append("svg:circle")
        .attr
          class: "percentile"
          r: 4
          cx: (d) -> x d[0]
          cy: (d) -> y d[1]

      viz.selectAll("line.fugi").data([0]).enter()
        .append("svg:line")
        .attr
          class: "fugi"
          y1: y(0)
          y2: y(1)
          x1: x
          x2: x

      viz.selectAll("circle.fugired").data([0]).enter()
        .append("svg:circle")
        .attr
          class: "fugired"
          r: 5
          cx: 0
          cy: 0

      viz.selectAll("circle.fugiblue").data([0]).enter()
        .append("svg:circle")
        .attr
          class: "fugiblue"
          r: 5
          cx: 0
          cy: 0

      $svg = $(viz[0][0]).closest("svg")
      $svg.mousemove((e) ->
        pos = x.invert(e.clientX - 10)
        if pos < axisTicks[0]
          pos = axisTicks[0]
        else pos = axisTicks[axisTicks.length - 1]  if pos > axisTicks[axisTicks.length - 1]
        $svg.find("line.fugi").attr("x1", x(pos)).attr "x2", x(pos)
        $svg.find(".fugired, .fugiblue").attr "cx", x(pos)
        $svg.find(".fugiblue").attr "cy", yk(kde(pos))
        $svg.find(".fugired").attr "cy", y(kde.inverseQuantile(pos))
      ).mouseover(->
        $svg.find("line.fugi").show()
      ).mouseout(->
        $svg.find("line.fugi").hide()
      )

    height = 400
    width = 600
    data = []
    $('textarea#data').val().split("\n").forEach((x) ->
      return if x.match(/^(#|\s*$)/)
      n = Number(x)
      # TODO warn about unparseable lines?
      data.push(n) unless isNaN(n)
    )

    data.sort(d3.ascending)

    cdf data

).call this
