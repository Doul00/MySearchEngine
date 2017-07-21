#! /usr/bin/env python3

import argparse
import pickle
import subprocess
import sys
from timeit import default_timer as timer
import os


import networkx as nx
import matplotlib.pyplot as plt
from matplotlib import cm
from mpl_toolkits.mplot3d import Axes3D


def main(argv):
    parser = argparse.ArgumentParser(description="Search Engine benchmarker.")
    parser.add_argument('--index', action='store', dest='index',
                        help='Max number of index.', type=int,
                        default=5)
    parser.add_argument('--go', action='store', dest='go',
                        help='Max number of goroutines.', type=int,
                        default=10)
    parser.add_argument('--mode', action='store', dest='mode',
                        help='Either "index" or "search"', required=True)
    parser.add_argument('--data', action='store', dest='data',
                        help='Data path.', required=True)

    args = parser.parse_args(argv)

    if args.mode == 'search':
        benchmark_search(args.index, args.go, args.data)
    elif args.mode == 'index':
        benchmark_index(args.index, args.data)
    else:
        raise Exception('Invalid mode: {}.'.format(args.mode))


def execute(args):
    begin = timer()
    subprocess.Popen(args, stdout=subprocess.DEVNULL,
                     stdin=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
                     shell=True).wait()
    return timer() - begin


def clean():
    execute('rm *.idx')


def build_index(data, nb):
    cmd = './myGoogle' +\
           ' index -go={} -index="tmp" -path="{}"'.format(nb, data)
    return execute(cmd)


def create_query():
    words = ['arthur', 'king', 'actress', 'obama', 'computer', 'apple',
             'chair', 'plane', 'wikipedia', 'name', 'is']
    words = list(map(lambda x: '"{}"'.format(x), words))
    echo = 'echo {} | '.format('\n'.join(words))
    cmd = echo + './myGoogle search -index="tmp_*.idx" -go={}'
    return cmd


def benchmark_index(nb_index, data):
    spaces = ' ' * 80 + '\r'
    msg = 'index: {} index\r'

    xs, ys = [], []
    clean()
    for idx in range(1, nb_index+1):
        print(spaces, end='')
        print(msg.format(idx), end='')

        t = build_index(data, str(idx))
        xs.append(idx)
        ys.append(t)

    print(spaces, end='')
    print('Finish! Computing figure...')

    fig = plt.figure()
    ax = fig.add_subplot(111)
    ax.plot(xs, ys)
    ax.set_xlabel('Number of indexes', labelpad=10)
    ax.set_ylabel('Time in seconds', labelpad=10)
    plt.show()


def benchmark_search(nb_index, nb_go, data):
    spaces = ' ' * 80 + '\r'
    msg = 'search: {} index, {} goroutines\r'
    query = create_query()

    xs, ys, zs = [], [], []
    for idx in range(1, nb_index+1):
        clean()
        build_index(data, str(idx))

        for go in range(1, nb_go+1):
            print(spaces, end='')
            print(msg.format(idx, go), end='')

            t = execute(query.format(str(go)))

            xs.append(idx)
            ys.append(go)
            zs.append(t)

    print(spaces, end='')
    print('Finish! Computing figure...')

    fig = plt.figure()
    ax = fig.add_subplot(111, projection='3d')
    ax.plot_trisurf(xs, ys, zs, cmap='jet')
    ax.set_xlabel('Number of indexes', labelpad=10)
    ax.set_ylabel('Number of workers', labelpad=10)
    ax.set_zlabel('Time in seconds', labelpad=10)
    plt.show()
    return xs, ys, zs


if __name__ == '__main__':
    main(sys.argv[1:])