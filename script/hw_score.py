# encoding=utf8

import xlrd
import sys
import psycopg2

wb = xlrd.open_workbook(sys.argv[1])
db = psycopg2.connect('dbname=kgo user=postgres')
cur=db.cursor()

def split_range(s):
    s=s.strip()
    if s.find('-')>=0:
        ranges=s.split('-')
        return round(float(ranges[0]),1),round(float(ranges[1]),1)
    elif s.find('<')>=0:
        ranges=s.split('<')
        return 0,round(float(ranges[1])-0.1,1)
    elif s.find('>')>=0:
        ranges=s.split('>')
        return round(float(ranges[1])+0.1,1),float('inf')

def read_sheet(sheet,gender):
    for i in range(sheet.nrows):
        ranges=split_range(sheet.cell_value(i,0))
        score1=split_range(sheet.cell_value(i,1))
        score2=split_range(sheet.cell_value(i,2))
        score3=split_range(sheet.cell_value(i,3))
        score4=split_range(sheet.cell_value(i,4))
        score5=split_range(sheet.cell_value(i,5))
        cur.execute('INSERT INTO "standard_scale_hw_score"("gender","height_min","height_max","weight_min","weight_max","score") VALUES(%s,%s,%s,%s,%s,%s)',(gender,ranges[0],ranges[1],score1[0],score1[1],1))
        cur.execute('INSERT INTO "standard_scale_hw_score"("gender","height_min","height_max","weight_min","weight_max","score") VALUES(%s,%s,%s,%s,%s,%s)',(gender,ranges[0],ranges[1],score2[0],score2[1],3))
        cur.execute('INSERT INTO "standard_scale_hw_score"("gender","height_min","height_max","weight_min","weight_max","score") VALUES(%s,%s,%s,%s,%s,%s)',(gender,ranges[0],ranges[1],score3[0],score3[1],5))
        cur.execute('INSERT INTO "standard_scale_hw_score"("gender","height_min","height_max","weight_min","weight_max","score") VALUES(%s,%s,%s,%s,%s,%s)',(gender,ranges[0],ranges[1],score4[0],score4[1],3))
        cur.execute('INSERT INTO "standard_scale_hw_score"("gender","height_min","height_max","weight_min","weight_max","score") VALUES(%s,%s,%s,%s,%s,%s)',(gender,ranges[0],ranges[1],score5[0],score5[1],1))
        # print("{0}: {1} - {2} | {3} - {4} | {5} - {6} | {7} - {8} | {9} - {10} | {11} - {12}".format(gender,ranges[0],ranges[1],score1[0],score1[1],score2[0],score2[1],score3[0],score3[1],score4[0],score4[1],score5[0],score5[1]))


read_sheet(wb.sheet_by_index(0),'male')
read_sheet(wb.sheet_by_index(1),'female')

db.commit()
cur.close()