import numpy as np
import pandas as pd 
import matplotlib.pyplot as plt
import seaborn as sns

# Jupiter notebook
# %matplotlib inline
# %reload_ext autoreload
# %autoreload 2

cs_df = pd.read_csv('sec_nanos.csv')
cs_df.head()

print(sns.get_dataset_names())