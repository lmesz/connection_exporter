# Connection exporter

* The aim of this small exporter is to collect socket statistics based on the outgoing port.

* The generated format is the following.
  ```
  connection_metric{port="58112",state="ESTABLISHED"} 1
  connection_metric{port="8080",state="TIME_WAIT"} 1
  ```
  It means there are 1 established connection towards 58112 which is in ESTABLISHED state and there is 1 connection towards 8080 which is in TIME_WAIT state.

# ToDo
[] tests
[] debug log
[] extract Collect to separate functions
[] better naming of variables
