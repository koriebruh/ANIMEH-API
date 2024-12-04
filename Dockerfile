FROM docker.elastic.co/elasticsearch/elasticsearch:8.10.0

# Pastikan plugin machine learning diaktifkan untuk mendukung operasi vektor
ENV xpack.ml.enabled=true
