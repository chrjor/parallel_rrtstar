#!/bin/bash

# Bash script that runs benchmark tests and calls python script of the same name
# to create speed-up graphs.
# 
# Example usage:
#     ./speedup.sh run
#     ./speedup.sh --help


# Display help information
Usage()
{
    echo "Usage: speedup.sh [option]" 
    echo
    echo "options:"
    echo "     run          run the benchmark"
    echo "     --help       display this help"
    echo 
}

# Check command line args passed to script
if [ $# -eq 0 ]; then
    Usage
    exit 1
elif [[ "$1" == "--help" ]]; then
    Usage
    exit 0
fi

# Create temporary output file
out="benchmark/output.txt"
rm $out
touch $out

# Write the column headers to file
header=("strategy" "difficulty" "threads" "test1" "test2" "test3" "test4" "test5")
echo "${header[@]}" >> $out

# Set the speedup test size(s)
tests=()
if [[ "$1" == "run" ]]; then
    tests=("4000 data/extraeasyMaze.txt" "8000 data/easyMaze.txt" "16000 data/mediumMaze.txt" "32000 data/hardMaze.txt")
else
    Usage
    exit 1
fi

# Execute sequential tests
for test in "${tests[@]}"; do
    level=$(echo "$test" | cut -d'/' -f2 | cut -d'M' -f1)
    sequential=( "s" $level "1" )
    for iter in {1..5}; do
        sequential+=($(go run proj3-redesigned/pathfinder bench $test))
    done
    echo "${sequential[@]}" >> $out
done

# Execute parallelized tests
for model in "ws" "bsp"; do
    for test in "${tests[@]}"; do
        level=$(echo "$test" | cut -d'/' -f2 | cut -d'M' -f1)
        for threads in 2 4 6 8 12; do
            threaded=( $model $level $threads )
            for iter in {1..5}; do 
                threaded+=($(go run proj3-redesigned/pathfinder bench $test $model $threads))
            done
            echo "${threaded[@]}" >> $out
        done
    done
done

# Create speedup graph and csv in python
python3 benchmark/speedup.py

# Clean up
rm $out
