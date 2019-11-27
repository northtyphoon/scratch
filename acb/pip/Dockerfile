FROM python:3.7

WORKDIR /root/wheels

RUN pip wheel --trusted-host pypi.python.org --trusted-host files.pythonhosted.org --trusted-host pypi.org --wheel-dir=/root/wheels numpy cython && \
    pip install --no-index --find-links=/root/wheels numpy cython