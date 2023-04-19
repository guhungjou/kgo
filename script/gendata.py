#!/usr/bin/python3
# encoding=utf8


import psycopg2
import random
import datetime

pg = psycopg2.connect(database='kgo', user='postgres', host='localhost')
cur = pg.cursor()

# year = random.randint(2020, 2021)
# month = random.randint(1, 12)
# day = random.randint(1, 28)
start = datetime.date(2021, 9, 1)

cur.execute('DELETE FROM kindergarten_student_morning_check')
cur.execute('DELETE FROM kindergarten_student_medical_examination')
cur.execute('DELETE FROM kindergarten_student_fitness_test')

cur.execute('SELECT id,kindergarten_id,gender,birthday FROM kindergarten_student')
students = cur.fetchall()


def getscore(name, gender, age, value):
    cur.execute('SELECT score FROM standard_scale_score WHERE name=%s AND gender=%s AND age=%s AND min<=%s AND max>=%s LIMIT 1',
                (name, gender, age, value, value))
    r = cur.fetchone()
    return r[0] if r else 0


def gethwscore(gender, height, weight):
    cur.execute('SELECT score FROM standard_scale_hw_score WHERE gender=%s AND height_min<=%s AND height_max>=%s AND weight_min<=%s AND weight_max>=%s LIMIT 1',
                (gender, height, height, weight, weight))
    r = cur.fetchone()
    return r[0] if r else 0


def getage(birthday):
    now = datetime.datetime.now()

    months = 12*(now.year-birthday.year) + \
        (int(now.month) - int(birthday.month))

    year = months / 12
    month = months % 12
    age = int(year)
    if month >= 4 and month <= 9:
        age += 0.5
    elif month > 9:
        age += 1
    return age


for r in students:
    d = start
    x = start
    y = start
    sid = r[0]
    kid = r[1]
    gender = r[2]
    birthday = r[3]
    age = getage(birthday)
    for i in range(220):
        temperature = round(random.normalvariate(36.6, 0.4), 2)
        status = 'normal'
        if temperature < 36:
            status = 'low'
        elif temperature > 37.3:
            status = 'high'
        cur.execute('INSERT INTO kindergarten_student_morning_check(student_id,kindergarten_id,temperature,date,created_at,updated_at,temperature_status) VALUES(%s,%s,%s,%s,%s,%s,%s)',
                    (sid, kid, temperature, d.isoformat(), d.isoformat(), d.isoformat(), status))
        d = d+datetime.timedelta(days=1)
    h = round(random.normalvariate(116, 4), 1)
    w = round(random.normalvariate(20.5, 2), 1)
    for i in range(31):
        hstatus = 'normal'
        wstatus = 'normal'
        bmi = round(w/h/h*10000, 2)
        bstatus = 'normal'
        if gender == 'male':
            if h < 111.2:
                hstatus = 'low'
            elif h > 120:
                hstatus = 'high'
            if w < 18.4:
                wstatus = 'low'
            elif w > 23.6:
                wstatus = 'high'
            if bmi < 13.7:
                bstatus = 'low'
            elif bmi > 17.2:
                bstatus = 'high'
        else:
            if h < 109.7:
                hstatus = 'low'
            elif h > 119.6:
                hstatus = 'high'
            if w < 17.3:
                wstatus = 'low'
            elif w > 22.9:
                wstatus = 'high'
            if bmi < 13.4:
                bstatus = 'low'
            elif bmi > 17.3:
                bstatus = 'high'
        date = x.isoformat()

        hb = round(random.normalvariate(135, 17), 2)
        hbstatus = 'normal'
        if hb < 110:
            hbstatus = 'low'
        elif hb > 160:
            hbstatus = 'high'

        alt = round(random.normalvariate(20, 16), 2)
        if alt < 0:
            alt = 0
        altstatus = 'normal'
        if alt > 40:
            altstatus = 'high'

        sightl = round(random.normalvariate(4.8, 0.2), 1)
        sightr = round(random.normalvariate(4.8, 0.2), 1)
        if sightl < 4:
            sightl = 4
        if sightr < 4:
            sightr = 4
        slstatus = 'normal'
        srstatus = 'normal'
        if sightl < 4.8:
            slstatus = 'low'
        if sightr < 4.8:
            srstatus = 'low'

        tooth = round(random.normalvariate(19, 2))
        caries = round(random.normalvariate(2, 2))
        if tooth < 0:
            tooth = 0
        if caries < 0:
            caries = 0

        cur.execute('INSERT INTO kindergarten_student_medical_examination(student_id,kindergarten_id,date,created_at,updated_at,height,weight,height_status,weight_status,height_updated_at,weight_updated_at,bmi,bmi_status,bmi_updated_at,hemoglobin,hemoglobin_status,hemoglobin_updated_at,alt,alt_status,alt_updated_at,sight_l,sight_r,sight_l_status,sight_r_status,sight_updated_at,tooth_count,tooth_caries_count,tooth_updated_at) VALUES(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)',
                    (sid, kid, date, date, date, h, w, hstatus, wstatus, date, date, bmi, bstatus, date, hb, hbstatus, date, alt, altstatus, date, sightl, sightr, slstatus, srstatus, date, tooth, caries, date))
        h = round(h + random.uniform(0, 0.15), 1)
        w = round(w+random.uniform(0, 0.1), 1)
        x = x+datetime.timedelta(weeks=1)

    shuttleRun10 = 16
    standingLongJump = 20
    baseballThrow = 1
    bunnyHopping = 25
    sitAndReach = 2.5
    balanceBeam = 35
    height = 100
    weight = 20
    for i in range(30):
        date = y.isoformat()
        height = round(height+random.uniform(0.8, 1.9), 0)
        weight = round(weight+random.uniform(0.3, 0.5), 1)

        h = height
        w = weight
        if gender == 'female':
            h -= 8
            w -= 3
        hwScore = gethwscore(gender, h, w)
        shuttleRun10 = round(shuttleRun10-random.uniform(0.2, 0.4), 1)
        shuttleRun10Score = getscore('10米折返跑(秒)', gender, age, shuttleRun10)
        standingLongJump = round(standingLongJump+random.uniform(1.5, 2.5), 1)
        standingLongJumpScore = getscore(
            '立定跳远(厘米)', gender, age, standingLongJump)
        baseballThrow = round(baseballThrow+random.choice([0, 0, 0, 0.5]), 1)
        baseballThrowScore = getscore('网球掷远(米)', gender, age, baseballThrow)
        bunnyHopping = round(bunnyHopping-random.uniform(0.5, 0.8), 1)
        bunnyHoppingScore = getscore('双脚连续跳(秒)', gender, age, bunnyHopping)
        sitAndReach = round(sitAndReach+random.uniform(0.2, 0.4), 1)
        sitAndReachScore = getscore('坐位体前屈(厘米)', gender, age, sitAndReach)
        balanceBeam = round(balanceBeam-random.uniform(0.5, 1.4), 1)
        balanceBeamScore = getscore('走平衡木(秒)', gender, age, balanceBeam)

        totalScore = hwScore+shuttleRun10Score+standingLongJumpScore + \
            baseballThrowScore+bunnyHoppingScore+sitAndReachScore+balanceBeamScore
        totalStatus = 'fail'
        if totalScore > 31:
            totalStatus = 'excellent'
        elif totalScore >= 28 and totalScore <= 31:
            totalStatus = 'good'
        elif totalScore >= 21 and totalScore <= 27:
            totalStatus = 'okay'
        cur.execute('INSERT INTO kindergarten_student_fitness_test(student_id,kindergarten_id,date,created_at,updated_at,shuttle_run_10,shuttle_run_10_score,shuttle_run_10_updated_at,standing_long_jump,standing_long_jump_score,standing_long_jump_updated_at,baseball_throw,baseball_throw_score,baseball_throw_updated_at,bunny_hopping,bunny_hopping_score,bunny_hopping_updated_at,sit_and_reach,sit_and_reach_score,sit_and_reach_updated_at,balance_beam,balance_beam_score,balance_beam_updated_at,total_score,height,height_updated_at,weight,weight_updated_at,height_and_weight_score,total_status) VALUES(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)',
                    (sid, kid, date, date, date, shuttleRun10, shuttleRun10Score, date, standingLongJump, standingLongJumpScore, date, baseballThrow, baseballThrowScore, date, bunnyHopping, bunnyHoppingScore, date, sitAndReach, sitAndReachScore, date, balanceBeam, balanceBeamScore, date, totalScore, h, date, w, date, hwScore, totalStatus))
        y = y+datetime.timedelta(weeks=1)


pg.commit()
cur.close()
pg.close()
