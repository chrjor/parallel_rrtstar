# Python script that processes benchmark data and produce speedup graph
# This script is called by process_speedup.sh, which runs the actual
# tests.


import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
import seaborn as sns

np.float = float

# Create dataframe
data = pd.read_table("benchmark/output.txt", sep='\s+')

# Calculate average for each test
data["avg"] = data.iloc[:,3:8].mean(axis=1)

# Separate out sequential and parallelized tests
num_tests = len(data["difficulty"].unique())
data_seq = data.iloc[:num_tests,[0,1,2,8]]
data_seq.set_index("difficulty", inplace=True)
data_par = data.iloc[num_tests:,[0,1,2,8]]

# Calculate speedup
def speedup(row):
    return data_seq.loc[row["difficulty"]]["avg"] / row["avg"]
data["speedup"] = data_par.apply(speedup, axis=1)

# Plot speedup graph
for model in data["strategy"].unique():
    if model != "s":
        data_plot = data[data["strategy"] == model]
        sns.set_style("whitegrid")
        sns.lineplot(x="threads", 
                    y="speedup", 
                    hue="difficulty", 
                    data=data_plot
                    ).set(
                        title=f"Robot Pathfinder Speedup ({model})", 
                        xlabel="Threads",
                        ylabel="Speedup")
        handles, labels = plt.gca().get_legend_handles_labels()
        new_labels = [label + f" ({num_samples} samples)" 
                      for label, num_samples in zip(labels, [0, 4000, 8000, 16000, 32000])]
        plt.legend(handles=handles[1:],
                   labels=new_labels[1:],
                   title="Maze Difficulty",
                   title_fontsize="small",
                   fontsize="x-small")

        # Create PNG and CSV
        plt.savefig(f"benchmark/speedup_{model}_graph.png")
        plt.clf()

data.to_csv("benchmark/speedup.csv")
