#!/bin/bash
#
#SBATCH --mail-user=christianj@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=pp-proj3
#SBATCH --output=benchmark/slurm/%j.stdout
#SBATCH --error=benchmark/slurm/%j.stderr
#SBATCH --chdir=/home/christianj/mpcs52060/autumn23/project-3-chrjor-1/proj3
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=4:00:00

module load golang/1.19
./benchmark/speedup.sh run
