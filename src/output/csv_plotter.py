"""
	Author: Max Miller
	Assignee: PHYS 605
	Purpose: Create a set of IV Curves from CSV's given by the AD2 ocislliscope
	Created: 2/16/2018
	Last Edit: 3/27/2018
"""

import glob
from matplotlib import pyplot as plt
import pandas as pd



def main():

	MANUAL_MODE = True # Set to True if you want to manually create the graphs using the command line

	extension = 'csv'
	result = [i for i in glob.glob('*.{}'.format(extension))]

	project_data = [] # Holds the dataframes generated from the CSV
	names_list = [] # Holds the names of the dataframes
	count = 0

	for item in result:
		if '' in item: # Fill in this string if you want a restriction on the types of CSV's you want to read in
			try:
				print('[{}]\t{}'.format(count, item))
				names_list.append(item)
				project_data.append(pd.read_csv(item, sep=",")) # removes all lines but the headers
				count+=1
			except pd.io.common.CParserError:
				print('CSV is not cleaned: remove headers.')

	print('Completed read')

	current_input = ''
	if MANUAL_MODE:
		while current_input is not 'exit':
			try:
				# Uncomment this for debugging. Allows you to pick a plot to graph
				print('Type \'exit\' to exit the program.')
				current_input = int(input('Choose a graph (use the number): ')) # Held over for debugging

				print('Selected ' + names_list[current_input])

				gen_x = raw_input('Set x-axis name: ') # Set this to be the x-axis name of the coorisponding column in your csv
				gen_y = raw_input('Set y-axis name: ') # Set this to be the y-axis name of the coorisponding column in your csv
				title = raw_input('Set a Title: ') 
				item = project_data[current_input]

				(item).plot(y=gen_y, x=gen_x, title=title) # Sets the x, y, and title (with the .csv extension removed)
				plt.savefig(names_list[current_input][0:-4] + '.png') # saves the figures to a new folder
			except KeyError:
				print('Error: Axis names are incorrect. Returning to top of program.')
			except NameError:
				print('Erorr: Please use only integers when selecting graphs. Example: for graph one, type 1. Returning to top of program.')

		print('Finished')
	else:
		count = 0 # Prooooobably shouldn't reuse the count here, but theres no harm
		gen_x = 'Second' # Set this to be the x-axis name of the coorisponding column in your csv
		gen_y = 'Capacitance Value' # Set this to be the y-axis name of the coorisponding column in your csv

		#Produces the plots
		for item in project_data:
			(item).plot(y=gen_y, x=gen_x, title=names_list[count][0:-4]) # Sets the x, y, and title (with the .csv extension removed)
			plt.savefig(names_list[count][0:-4] + '.png') # saves the figures to a new folder
			count+=1

	# plt.show() # Use to show the graph to the user
	
if __name__ == '__main__':
        main()
