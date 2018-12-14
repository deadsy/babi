#!/usr/bin/env python3

import math
import plotly

STEPS = 1000
yvals = [0,] * STEPS
xvals = [0,] * STEPS

def mapExp(x, y0, y1, k):
	k0 = 1.0 / (1.0 - math.pow(2.0, k))
	a = (y0 - y1) * k0
	b = y0 - a
	return (a * math.pow(2.0, k*x)) + b

def main():

  k = 6
  y0 = 2.0
  y1 = -3.0

  for i in range(STEPS):
    xvals[i] = i / STEPS
    yvals[i] = mapExp(xvals[i], y0, y1, k)


  data = [
    plotly.graph_objs.Scatter(
      x=xvals,
      y=yvals,
      mode = 'lines',
    ),
  ]
  layout = plotly.graph_objs.Layout(
    title='Exponential Mapping',
    xaxis=dict(
      title='x',
    ),
    yaxis=dict(
      title='y',
      rangemode='tozero',
    ),
  )

  figure = plotly.graph_objs.Figure(data=data, layout=layout)
  plotly.offline.plot(figure, filename='map.html')

main()
