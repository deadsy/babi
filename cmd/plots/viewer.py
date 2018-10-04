#!/usr/bin/python3

import plotly
import plotly.graph_objs as go
import babi

data = [
  go.Scatter(
    x=babi.time,
    y=babi.amplitude,
    #mode = 'markers',
    mode = 'lines',
  ),
]

layout = go.Layout(
  title=babi.title,
  xaxis=dict(
    title="time",
  ),
  yaxis=dict(
    title="amplitude",
    rangemode="tozero",
  ),
)

figure = go.Figure(data=data, layout=layout)
plotly.offline.plot(figure, filename='babi.html')
