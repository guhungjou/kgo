CREATE DATABASE "kgo" WITH ENCODING=UTF8;

\c "kgo";

CREATE EXTENSION "uuid-ossp";
CREATE EXTENSION "pg_trgm";


CREATE TABLE "admin_user"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "username" VARCHAR(32) NOT NULL,
    "name" VARCHAR(32) NOT NULL,
    "salt" VARCHAR(32) NOT NULL,
    "ptype" VARCHAR(16) NOT NULL,
    "password" VARCHAR(256) NOT NULL,
    "status" VARCHAR(32) NOT NULL DEFAULT 'OK',
    "is_superuser" BOOLEAN NOT NULL DEFAULT FALSE,
    "phone" VARCHAR(32) NOT NULL,
    "created_by" BIGINT,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE("username")
);
CREATE INDEX ON "admin_user"("phone");
CREATE INDEX ON "admin_user"("status");

CREATE TABLE "admin_token"(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_id" BIGINT NOT NULL,
    "device" VARCHAR(1024) NOT NULL,
    "ip" VARCHAR(256) NOT NULL,
    "expires_at" TIMESTAMP WITH TIME ZONE,
    "status" VARCHAR(32) NOT NULL DEFAULT 'OK',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON "admin_token"("user_id","status");
CREATE INDEX ON "admin_token"("user_id","expires_at");
CREATE INDEX ON "admin_token"("created_at");


/* 幼儿园 */
CREATE TABLE "kindergarten"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,        /* 主键 */
    "name" VARCHAR(128) NOT NULL,               /* 幼儿园名 */
    "remark" VARCHAR(256) NOT NULL DEFAULT '',  /* 备注信息 */
    "number_of_student" BIGINT NOT NULL DEFAULT 0, /* 学生数量 */
    "number_of_teacher" BIGINT NOT NULL DEFAULT 0, /* 老师数量 */
    "district_id" VARCHAR(32) NOT NULL DEFAULT '', /* 所在省市区 */
    "created_by" BIGINT,                           /* 创建人 */
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,    /* 创建时间 */
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,    /* 更新时间 */
    "deleted" BOOLEAN NOT NULL DEFAULT FALSE                        /* 是否删除 */
);
CREATE INDEX ON "kindergarten" USING gin("name" gin_trgm_ops);
CREATE INDEX ON "kindergarten" USING gin("remark" gin_trgm_ops);
CREATE INDEX ON "kindergarten"("district_id");
CREATE INDEX ON "kindergarten"("deleted");


/* 幼儿园班级 */
CREATE TABLE "kindergarten_class"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,/* 主键 */
    "name" VARCHAR(128) NOT NULL,   /* 班级名 */
    "remark" VARCHAR(256) NOT NULL DEFAULT '',  /* 备注信息 */
    "kindergarten_id" BIGINT NOT NULL,          /* 幼儿园ID */
    "number_of_student" BIGINT NOT NULL DEFAULT 0, /* 学生数量 */
    "number_of_teacher" BIGINT NOT NULL DEFAULT 0, /* 老师数量 */
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted" BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX ON "kindergarten_class" USING gin("name" gin_trgm_ops);
CREATE INDEX ON "kindergarten_class" USING gin("remark" gin_trgm_ops);
CREATE INDEX ON "kindergarten_class"("kindergarten_id");
CREATE INDEX ON "kindergarten_class"("deleted");

/* 幼儿园老师（包括园长） */
CREATE TABLE "kindergarten_teacher"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,    /* 主键 */
    "username" VARCHAR(32) NOT NULL,    /* 登录用户名 */
    "name" VARCHAR(32) NOT NULL,        /* 姓名 */
    "role" VARCHAR(32) NOT NULL,        /* 角色:园长，老师 */
    "gender" VARCHAR(32) NOT NULL,      /* 性别 m男，f女 */
    "phone" VARCHAR(32) NOT NULL DEFAULT '', /* 联系电话 */

    "kindergarten_id" BIGINT NOT NULL,      /* 幼儿园，不能为空，创建时确定 */
    "class_id" BIGINT NOT NULL DEFAULT 0, /* 班级 */

    "salt" VARCHAR(32) NOT NULL,
    "ptype" VARCHAR(16) NOT NULL,
    "password" VARCHAR(256) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted" BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE("username")
);
CREATE INDEX ON "kindergarten_teacher" USING gin("name" gin_trgm_ops);
CREATE INDEX ON "kindergarten_teacher" USING gin("phone" gin_trgm_ops);
CREATE INDEX ON "kindergarten_teacher"("kindergarten_id");
CREATE INDEX ON "kindergarten_teacher"("kindergarten_id","role");
CREATE INDEX ON "kindergarten_teacher"("role");
CREATE INDEX ON "kindergarten_teacher"("deleted");


CREATE TABLE "kindergarten_teacher_token"(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "teacher_id" BIGINT NOT NULL,
    "device" VARCHAR(1024) NOT NULL,
    "ip" VARCHAR(256) NOT NULL,
    "expires_at" TIMESTAMP WITH TIME ZONE,
    "status" VARCHAR(32) NOT NULL DEFAULT 'OK',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON "kindergarten_teacher_token"("teacher_id","status");
CREATE INDEX ON "kindergarten_teacher_token"("teacher_id","expires_at");
CREATE INDEX ON "kindergarten_teacher_token"("created_at");

/* 学生 */
CREATE TABLE "kindergarten_student"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "kindergarten_id" BIGINT NOT NULL,  /* 幼儿园ID */
    "no" VARCHAR(64) NOT NULL DEFAULT '', /* 学号 */
    "class_id" BIGINT NOT NULL,         /* 班级ID */
    "name" VARCHAR(64) NOT NULL,        /* 学生名 */
    "gender" VARCHAR(32) NOT NULL,      /* 性别，male男，female女 */
    "birthday" DATE NOT NULL,           /* 生日 */
    "device" VARCHAR(256) NOT NULL,     /* 手环MAC地址 */
    "remark" VARCHAR(256) NOT NULL DEFAULT '', /* 备注信息 */
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted" BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX ON "kindergarten_student" USING gin("name" gin_trgm_ops);
CREATE INDEX ON "kindergarten_student"("kindergarten_id");
CREATE INDEX ON "kindergarten_student"("class_id");
CREATE INDEX ON "kindergarten_student"("gender");
CREATE INDEX ON "kindergarten_student"("device");
CREATE INDEX ON "kindergarten_student"("deleted");
CREATE INDEX ON "kindergarten_student"("no");


/* 晨检 */
CREATE TABLE "kindergarten_student_morning_check"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "kindergarten_id" BIGINT NOT NULL,
    "student_id" BIGINT NOT NULL,
    "date" DATE NOT NULL,
    "temperature" FLOAT NOT NULL DEFAULT 0,
    "temperature_status" VARCHAR(32) NOT NULL DEFAULT '',
    "hand" VARCHAR(256) NOT NULL DEFAULT '',
    "mouth" VARCHAR(256) NOT NULL DEFAULT '',
    "eye" VARCHAR(256) NOT NULL DEFAULT '',
    "hand_status" VARCHAR(32) NOT NULL DEFAULT '',
    "mouth_status" VARCHAR(32) NOT NULL DEFAULT '',
    "eye_status" VARCHAR(32) NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted" BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX ON "kindergarten_student_morning_check"("kindergarten_id");
CREATE INDEX ON "kindergarten_student_morning_check"("student_id");
CREATE INDEX ON "kindergarten_student_morning_check"("date");
CREATE INDEX ON "kindergarten_student_morning_check"("temperature");
CREATE INDEX ON "kindergarten_student_morning_check"("deleted");
CREATE INDEX ON "kindergarten_student_morning_check"("temperature_status");
CREATE INDEX ON "kindergarten_student_morning_check"("hand_status");
CREATE INDEX ON "kindergarten_student_morning_check"("mouth_status");
CREATE INDEX ON "kindergarten_student_morning_check"("eye_status");

/* 体检 */
CREATE TABLE "kindergarten_student_medical_examination"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "kindergarten_id" BIGINT NOT NULL,
    "student_id" BIGINT NOT NULL,
    "date" DATE NOT NULL,
    "height" FLOAT NOT NULL DEFAULT 0, /* 身高 CM */
    "height_updated_at" TIMESTAMP WITH TIME ZONE,
    "height_status" VARCHAR(32) NOT NULL DEFAULT '',

    "weight" FLOAT NOT NULL DEFAULT 0, /* 体重 KG */
    "weight_updated_at" TIMESTAMP WITH TIME ZONE,
    "weight_status" VARCHAR(32) NOT NULL DEFAULT '',

    "bmi" FLOAT NOT NULL DEFAULT 0, /* 体重比身高 */
    "bmi_updated_at" TIMESTAMP WITH TIME ZONE,
    "bmi_status" VARCHAR(32) NOT NULL DEFAULT '',

    "hemoglobin" FLOAT NOT NULL DEFAULT 0, /* 血红蛋白 */
    "hemoglobin_remark" VARCHAR(256) NOT NULL DEFAULT '', /* 血红蛋白备注 */
    "hemoglobin_status" VARCHAR(32) NOT NULL DEFAULT '',
    "hemoglobin_updated_at" TIMESTAMP WITH TIME ZONE,

    "tooth_count" INT NOT NULL DEFAULT 0, /* 牙齿数 */
    "tooth_caries_count" INT NOT NULL DEFAULT 0, /* 蛀牙数 */
    "tooth_remark" VARCHAR(256) NOT NULL DEFAULT '', /* 牙齿备注 */
    "tooth_updated_at" TIMESTAMP WITH TIME ZONE,

    -- "sight_l" FLOAT NOT NULL DEFAULT 0, /* 左眼视力 */
    "sight_l_s" VARCHAR(256) NOT NULL DEFAULT '', -- 左眼球镜度
    "sight_l_c" VARCHAR(256) NOT NULL DEFAULT '', -- 左眼柱镜度
    "sight_l_remark" VARCHAR(256) NOT NULL DEFAULT '',/* 左眼备注 */
    "sight_l_s_status" VARCHAR(32) NOT NULL DEFAULT '', -- 左眼球镜状态
    "sight_l_c_status" VARCHAR(32) NOT NULL DEFAULT '', -- 左眼柱镜状态
    -- "sight_l_status" VARCHAR(32) NOT NULL DEFAULT '',
    -- "sight_r" FLOAT NOT NULL DEFAULT 0, /* 右眼视力 */
    "sight_r_s" VARCHAR(256) NOT NULL DEFAULT '', -- 右眼球镜度
    "sight_r_c" VARCHAR(256) NOT NULL DEFAULT '', -- 右眼柱镜度
    "sight_r_remark" VARCHAR(256) NOT NULL DEFAULT '',/* 右眼备注 */
    "sight_r_s_status" VARCHAR(32) NOT NULL DEFAULT '', -- 右眼球镜状态
    "sight_r_c_status" VARCHAR(32) NOT NULL DEFAULT '', -- 右眼柱镜状态
    -- "sight_r_status" VARCHAR(32) NOT NULL DEFAULT '',
    "sight_updated_at" TIMESTAMP WITH TIME ZONE,

    "alt" FLOAT NOT NULL DEFAULT 0, /* 谷丙转氨酶 */
    "alt_remark" VARCHAR(256) NOT NULL DEFAULT '', /* 谷丙转氨酶 备注 */
    "alt_status" VARCHAR(32) NOT NULL DEFAULT '',
    "alt_updated_at" TIMESTAMP WITH TIME ZONE,
    
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted" BOOLEAN NOT NULL DEFAULT FALSE,
    "eye_left" FLOAT NOT NULL DEFAULT 0,
    "eye_right" FLOAT NOT NULL DEFAULT 0
);
CREATE INDEX ON "kindergarten_student_medical_examination"("kindergarten_id");
CREATE INDEX ON "kindergarten_student_medical_examination"("student_id");
CREATE INDEX ON "kindergarten_student_medical_examination"("date");
CREATE INDEX ON "kindergarten_student_medical_examination"("deleted");
CREATE INDEX ON "kindergarten_student_medical_examination"("height_status");
CREATE INDEX ON "kindergarten_student_medical_examination"("weight_status");
CREATE INDEX ON "kindergarten_student_medical_examination"("hemoglobin_status");
CREATE INDEX ON "kindergarten_student_medical_examination"("alt_status");
CREATE INDEX ON "kindergarten_student_medical_examination"("sight_l_status");
CREATE INDEX ON "kindergarten_student_medical_examination"("sight_r_status");
CREATE INDEX ON "kindergarten_student_medical_examination"("bmi_status");
CREATE INDEX ON "kindergarten_student_medical_examination"("updated_at");



/* 体测 */
CREATE TABLE "kindergarten_student_fitness_test"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "kindergarten_id" BIGINT NOT NULL,
    "student_id" BIGINT NOT NULL,
    "date" DATE NOT NULL,

    "height" FLOAT NOT NULL DEFAULT 0, -- 身高
    "height_updated_at" TIMESTAMP WITH TIME ZONE,

    "weight" FLOAT NOT NULL DEFAULT 0, -- 体重
    "weight_updated_at" TIMESTAMP WITH TIME ZONE,

    "height_and_weight_score" FLOAT NOT NULL DEFAULT 0, -- 身高体重评分

    "shuttle_run_10" FLOAT NOT NULL DEFAULT 0,  -- 10米折返跑(秒)
    "shuttle_run_10_score" FLOAT NOT NULL DEFAULT 0,
    "shuttle_run_10_updated_at" TIMESTAMP WITH TIME ZONE,

    "standing_long_jump" FLOAT NOT NULL DEFAULT 0,  -- 立定跳远(厘米)
    "standing_long_jump_score" FLOAT NOT NULL DEFAULT 0,
    "standing_long_jump_updated_at" TIMESTAMP WITH TIME ZONE,

    "baseball_throw" FLOAT NOT NULL DEFAULT 0, -- 网球掷远(米)
    "baseball_throw_score" FLOAT NOT NULL DEFAULT 0,
    "baseball_throw_updated_at" TIMESTAMP WITH TIME ZONE,

    "bunny_hopping" FLOAT NOT NULL DEFAULT 0,  -- 双脚连续跳(秒)
    "bunny_hopping_score" FLOAT NOT NULL DEFAULT 0,
    "bunny_hopping_updated_at" TIMESTAMP WITH TIME ZONE,

    "sit_and_reach" FLOAT NOT NULL DEFAULT 0,   -- 坐位体前屈(厘米)
    "sit_and_reach_score" FLOAT NOT NULL DEFAULT 0,
    "sit_and_reach_updated_at" TIMESTAMP WITH TIME ZONE,

    "balance_beam" FLOAT NOT NULL DEFAULT 0,    -- 走平衡木(秒)
    "balance_beam_score" FLOAT NOT NULL DEFAULT 0,
    "balance_beam_updated_at" TIMESTAMP WITH TIME ZONE,

    "total_score" FLOAT NOT NULL DEFAULT 0, -- 总分
    "total_status" VARCHAR(32) NOT NULL DEFAULT '', -- 总评价

    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted" BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX ON "kindergarten_student_fitness_test"("kindergarten_id");
CREATE INDEX ON "kindergarten_student_fitness_test"("student_id");
CREATE INDEX ON "kindergarten_student_fitness_test"("date");
CREATE INDEX ON "kindergarten_student_fitness_test"("deleted");
CREATE INDEX ON "kindergarten_student_fitness_test"("shuttle_run_10_updated_at");
CREATE INDEX ON "kindergarten_student_fitness_test"("shuttle_run_10_score");
CREATE INDEX ON "kindergarten_student_fitness_test"("baseball_throw_updated_at");
CREATE INDEX ON "kindergarten_student_fitness_test"("baseball_throw_score");
CREATE INDEX ON "kindergarten_student_fitness_test"("standing_long_jump_updated_at");
CREATE INDEX ON "kindergarten_student_fitness_test"("standing_long_jump_score");
CREATE INDEX ON "kindergarten_student_fitness_test"("bunny_hopping_updated_at");
CREATE INDEX ON "kindergarten_student_fitness_test"("bunny_hopping_score");
CREATE INDEX ON "kindergarten_student_fitness_test"("sit_and_reach_updated_at");
CREATE INDEX ON "kindergarten_student_fitness_test"("sit_and_reach_score");
CREATE INDEX ON "kindergarten_student_fitness_test"("balance_beam_updated_at");
CREATE INDEX ON "kindergarten_student_fitness_test"("balance_beam_score");
CREATE INDEX ON "kindergarten_student_fitness_test"("updated_at");
CREATE INDEX ON "kindergarten_student_fitness_test"("total_status");
CREATE INDEX ON "kindergarten_student_fitness_test"("height_updated_at");
CREATE INDEX ON "kindergarten_student_fitness_test"("height_and_weight_score");
CREATE INDEX ON "kindergarten_student_fitness_test"("weight_updated_at");


CREATE TABLE "standard_scale_hl"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(32) NOT NULL, -- 类型
    "gender" VARCHAR(32) NOT NULL, -- 性别
    "age" FLOAT NOT NULL, -- 年龄
    "min" FLOAT NOT NULL, -- 下限
    "max" FLOAT NOT NULL, -- 上限,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE("name","gender","age")
);
CREATE INDEX ON "standard_scale_hl"("name");
CREATE INDEX ON "standard_scale_hl"("age");
CREATE INDEX ON "standard_scale_hl"("gender");

/* https://www.youlai.cn/dise/imagedetail/157_32867.html */
INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('身高', 2, 'male', 84.3, 91),
    ('身高', 2.5, 'male',88.9, 98.5),
    ('身高', 3, 'male', 91.1, 98.7),
    ('身高', 3.5, 'male', 95, 103.1),
    ('身高', 4, 'male', 98.7, 107.2),
    ('身高', 4.5, 'male', 102.1, 111),
    ('身高', 5, 'male', 105.3, 114.5),
    ('身高', 5.5, 'male', 108.4, 117.8),
    ('身高', 6, 'male', 111.2, 121.0),
    ('身高', 7, 'male', 116.6, 126.8),
    ('身高', 8, 'male', 121.6, 132.2),
    ('身高', 9, 'male', 126.5, 137.8),
    ('身高', 10, 'male', 131.4, 143.6);

INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('身高', 2, 'female', 83.8, 89.8),
    ('身高', 2.5, 'female',87.9, 94.7),
    ('身高', 3, 'female', 90.2, 98.1),
    ('身高', 3.5, 'female', 94, 101.8),
    ('身高', 4, 'female', 97.6, 105.7),
    ('身高', 4.5, 'female', 100.9, 109.3),
    ('身高', 5, 'female', 104.0, 112.8),
    ('身高', 5.5, 'female', 106.9, 116.2),
    ('身高', 6, 'female', 109.7, 119.6),
    ('身高', 7, 'female', 115.1, 126.2),
    ('身高', 8, 'female', 120.4, 132.4),
    ('身高', 9, 'female', 125.7, 138.7),
    ('身高', 10, 'female', 131.5, 145.1);

INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('体重', 2, 'male', 11.2, 14.0),
    ('体重', 2.5, 'male', 12.1, 15.3),
    ('体重', 3, 'male', 13.0, 16.4),
    ('体重', 3.5, 'male', 13.9, 17.6),
    ('体重', 4, 'male', 14.8, 18.7),
    ('体重', 4.5, 'male', 15.7, 19.9),
    ('体重', 5, 'male', 16.6, 21.1),
    ('体重', 5.5, 'male', 17.4, 22.3),
    ('体重', 6, 'male', 18.4, 23.6),
    ('体重', 7, 'male', 20.2, 26.5),
    ('体重', 8, 'male', 22.2, 30.0),
    ('体重', 9, 'male', 24.3, 34.0),
    ('体重', 10, 'male', 26.8, 38.7);

INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('体重', 2, 'female', 10.6, 13.2),
    ('体重', 2.5, 'female', 11.7, 14.7),
    ('体重', 3, 'female', 12.6, 16.1),
    ('体重', 3.5, 'female', 13.5, 17.2),
    ('体重', 4, 'female', 14.3, 18.3),
    ('体重', 4.5, 'female', 15, 19.4),
    ('体重', 5, 'female', 15.7, 20.4),
    ('体重', 5.5, 'female', 16.5, 21.6),
    ('体重', 6, 'female', 17.3, 22.9),
    ('体重', 7, 'female', 19.1, 26),
    ('体重', 8, 'female', 21.4, 30.2),
    ('体重', 9, 'female', 24.1, 35.3),
    ('体重', 10, 'female', 27.2, 40.9);


INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('血红蛋白', 0, '', 110, 160);

INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('谷丙转氨酶', 0, '', 0, 40);

INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('视力', 0, '', 4.8, -1);

INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('球镜度', 0, '', -1, 1);

INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('柱镜度', 0, '', -0.5, 0.5);

INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('体温', 0, '', 36, 37.3);

-- https://www.cdc.gov/growthcharts/data/set1clinical/cj41l023.pdf
INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('BMI', 2, 'male', 14.8, 18.2),
    ('BMI', 2.5, 'male', 14.6, 17.8),
    ('BMI', 3, 'male', 14.4, 17.4),
    ('BMI', 3.5, 'male', 14.2, 17.1),
    ('BMI', 4, 'male', 14, 17),
    ('BMI', 4.5, 'male', 14, 16.8),
    ('BMI', 5, 'male', 13.8, 16.8),
    ('BMI', 5.5, 'male', 13.8, 16.9),
    ('BMI', 6, 'male', 13.8, 17),
    ('BMI', 6.5, 'male', 13.7, 17.2),
    ('BMI', 7, 'male', 13.8, 17.4),
    ('BMI', 7.5, 'male', 13.8, 17.6),
    ('BMI', 8, 'male', 13.8, 17.9),
    ('BMI', 8.5, 'male', 13.8, 18.3),
    ('BMI', 9, 'male', 14, 18.6),
    ('BMI', 9.5, 'male', 14.1, 18.9),
    ('BMI', 10, 'male', 14.2, 19.4);

-- https://www.cdc.gov/growthcharts/data/set1clinical/cj41l024.pdf
INSERT INTO "standard_scale_hl"("name","age","gender","min","max") VALUES
    ('BMI', 2, 'female', 14.4, 18),
    ('BMI', 2.5, 'female', 14.2, 17.6),
    ('BMI', 3, 'female', 14, 17.2),
    ('BMI', 3.5, 'female', 13.8, 17),
    ('BMI', 4, 'female', 13.8, 16.8),
    ('BMI', 4.5, 'female', 13.6, 16.8),
    ('BMI', 5, 'female', 13.5, 16.8),
    ('BMI', 5.5, 'female', 13.5, 16.9),
    ('BMI', 6, 'female', 13.4, 17.1),
    ('BMI', 6.5, 'female', 13.4, 17.3),
    ('BMI', 7, 'female', 13.4, 17.6),
    ('BMI', 7.5, 'female', 13.5, 17.9),
    ('BMI', 8, 'female', 13.6, 18.3),
    ('BMI', 8.5, 'female', 13.6, 18.5),
    ('BMI', 9, 'female', 13.8, 19.1),
    ('BMI', 9.5, 'female', 13.9, 19.3),
    ('BMI', 10, 'female', 14, 19.7);

CREATE TABLE "standard_scale_score"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(32) NOT NULL, -- 类型
    "gender" VARCHAR(32) NOT NULL, -- 性别
    "age" FLOAT NOT NULL, -- 年龄
    "min" FLOAT NOT NULL, -- 下限
    "max" FLOAT NOT NULL, -- 上限,
    "score" FLOAT NOT NULL, -- 分数
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE("name","gender","age","min","max")
);
CREATE INDEX ON "standard_scale_score"("name");
CREATE INDEX ON "standard_scale_score"("age");
CREATE INDEX ON "standard_scale_score"("gender");
CREATE INDEX ON "standard_scale_score"("min");
CREATE INDEX ON "standard_scale_score"("max");

INSERT INTO "standard_scale_score"("name","gender","age","min","max","score") VALUES
    ('10米折返跑(秒)','male', 3, 12.9, 15.8, 1),
    ('10米折返跑(秒)','male', 3, 10.1, 12.8, 2),
    ('10米折返跑(秒)','male', 3, 9.1, 10.2, 3),
    ('10米折返跑(秒)','male', 3, 8.0, 9.0, 4),
    ('10米折返跑(秒)','male', 3, 0, 7.9, 5),
    ('10米折返跑(秒)','female', 3, 13.5, 16.8, 1),
    ('10米折返跑(秒)','female', 3, 10.6, 13.4, 2),
    ('10米折返跑(秒)','female', 3, 9.4, 10.5, 3),
    ('10米折返跑(秒)','female', 3, 8.2, 9.3, 4),
    ('10米折返跑(秒)','female', 3, 0, 8.1, 5),
    ('10米折返跑(秒)','male', 3.5, 11.4, 14, 1),
    ('10米折返跑(秒)','male', 3.5, 9.5, 11.3, 2),
    ('10米折返跑(秒)','male', 3.5, 8.4, 9.4, 3),
    ('10米折返跑(秒)','male', 3.5, 7.5, 8.3, 4),
    ('10米折返跑(秒)','male', 3.5, 0, 7.4, 5),
    ('10米折返跑(秒)','female', 3.5, 12.1, 14.9, 1),
    ('10米折返跑(秒)','female', 3.5, 9.8, 12, 2),
    ('10米折返跑(秒)','female', 3.5, 8.7, 9.7, 3),
    ('10米折返跑(秒)','female', 3.5, 7.7, 8.6, 4),
    ('10米折返跑(秒)','female', 3.5, 0, 7.6, 5),
    ('10米折返跑(秒)','male', 4, 10.2, 12.4, 1),
    ('10米折返跑(秒)','male', 4, 8.6, 10.1, 2),
    ('10米折返跑(秒)','male', 4, 7.7, 8.5, 3),
    ('10米折返跑(秒)','male', 4, 6.9, 7.6, 4),
    ('10米折返跑(秒)','male', 4, 0, 6.8, 5),
    ('10米折返跑(秒)','female', 4, 10.7, 13.2, 1),
    ('10米折返跑(秒)','female', 4, 9.1, 10.8, 2),
    ('10米折返跑(秒)','female', 4, 8.1, 9.0, 3),
    ('10米折返跑(秒)','female', 4, 7.2, 8.0, 4),
    ('10米折返跑(秒)','female', 4, 0, 7.1, 5),
    ('10米折返跑(秒)','male', 4.5, 9.8, 11.8, 1),
    ('10米折返跑(秒)','male', 4.5, 8.1, 9.7, 2),
    ('10米折返跑(秒)','male', 4.5, 7.3, 8.0, 3),
    ('10米折返跑(秒)','male', 4.5, 6.7, 7.2, 4),
    ('10米折返跑(秒)','male', 4.5, 0, 6.6, 5),
    ('10米折返跑(秒)','female', 4.5, 10.3, 12.4, 1),
    ('10米折返跑(秒)','female', 4.5, 8.6, 10.2, 2),
    ('10米折返跑(秒)','female', 4.5, 7.7, 8.5, 3),
    ('10米折返跑(秒)','female', 4.5, 7.0, 7.6, 4),
    ('10米折返跑(秒)','female', 4.5, 0, 6.9, 5),
    ('10米折返跑(秒)','male', 5, 9, 10.3, 1),
    ('10米折返跑(秒)','male', 5, 7.7, 8.9, 2),
    ('10米折返跑(秒)','male', 5, 7.0, 7.6, 3),
    ('10米折返跑(秒)','male', 5, 6.4, 6.9, 4),
    ('10米折返跑(秒)','male', 5, 0, 6.3, 5),
    ('10米折返跑(秒)','female', 5, 9.7, 11.2, 1),
    ('10米折返跑(秒)','female', 5, 8.1, 9.6, 2),
    ('10米折返跑(秒)','female', 5, 7.3, 8.0, 3),
    ('10米折返跑(秒)','female', 5, 6.7, 7.2, 4),
    ('10米折返跑(秒)','female', 5, 0, 6.6, 5),
    ('10米折返跑(秒)','male', 5.5, 8.6, 10, 1),
    ('10米折返跑(秒)','male', 5.5, 7.4, 8.5, 2),
    ('10米折返跑(秒)','male', 5.5, 6.8, 7.3, 3),
    ('10米折返跑(秒)','male', 5.5, 6.2, 6.7, 4),
    ('10米折返跑(秒)','male', 5.5, 0, 6.1, 5),
    ('10米折返跑(秒)','female', 5.5, 9.1, 10.2, 1),
    ('10米折返跑(秒)','female', 5.5, 7.7, 9, 2),
    ('10米折返跑(秒)','female', 5.5, 7.0, 7.6, 3),
    ('10米折返跑(秒)','female', 5.5, 6.4, 6.9, 4),
    ('10米折返跑(秒)','female', 5.5, 0, 6.3, 5),
    ('10米折返跑(秒)','male', 6, 8, 9.4, 1),
    ('10米折返跑(秒)','male', 6, 6.8, 7.9, 2),
    ('10米折返跑(秒)','male', 6, 6.3, 6.7, 3),
    ('10米折返跑(秒)','male', 6, 5.8, 6.2, 4),
    ('10米折返跑(秒)','male', 6, 0, 5.7, 5),
    ('10米折返跑(秒)','female', 6, 8.6, 10.2, 1),
    ('10米折返跑(秒)','female', 6, 7.3, 8.5, 2),
    ('10米折返跑(秒)','female', 6, 6.6, 7.2, 3),
    ('10米折返跑(秒)','female', 6, 6.1, 6.5, 4),
    ('10米折返跑(秒)','female', 6, 0, 6, 5);

INSERT INTO "standard_scale_score"("name","gender","age","min","max","score") VALUES
    ('立定跳远(厘米)','male', 3, 21, 29, 1),
    ('立定跳远(厘米)','male', 3, 30, 43, 2),
    ('立定跳远(厘米)','male', 3, 44, 58, 3),
    ('立定跳远(厘米)','male', 3, 59, 76, 4),
    ('立定跳远(厘米)','male', 3, 77, 'Infinity', 5),
    ('立定跳远(厘米)','female', 3, 21, 28, 1),
    ('立定跳远(厘米)','female', 3, 29, 39, 2),
    ('立定跳远(厘米)','female', 3, 40, 54, 3),
    ('立定跳远(厘米)','female', 3, 55, 71, 4),
    ('立定跳远(厘米)','female', 3, 72, 'Infinity', 5),
    ('立定跳远(厘米)','male', 3.5, 27, 34, 1),
    ('立定跳远(厘米)','male', 3.5, 35, 52, 2),
    ('立定跳远(厘米)','male', 3.5, 53, 69, 3),
    ('立定跳远(厘米)','male', 3.5, 70, 84, 4),
    ('立定跳远(厘米)','male', 3.5, 85, 'Infinity', 5),
    ('立定跳远(厘米)','female', 3.5, 25, 33, 1),
    ('立定跳远(厘米)','female', 3.5, 34, 49, 2),
    ('立定跳远(厘米)','female', 3.5, 50, 64, 3),
    ('立定跳远(厘米)','female', 3.5, 65, 81, 4),
    ('立定跳远(厘米)','female', 3.5, 82, 'Infinity', 5),
    ('立定跳远(厘米)','male', 4, 35, 46, 1),
    ('立定跳远(厘米)','male', 4, 47, 65, 2),
    ('立定跳远(厘米)','male', 4, 66, 79, 3),
    ('立定跳远(厘米)','male', 4, 80, 95, 4),
    ('立定跳远(厘米)','male', 4, 96, 'Infinity', 5),
    ('立定跳远(厘米)','female', 4, 32, 43, 1),
    ('立定跳远(厘米)','female', 4, 44, 59, 2),
    ('立定跳远(厘米)','female', 4, 60, 73, 3),
    ('立定跳远(厘米)','female', 4, 74, 89, 4),
    ('立定跳远(厘米)','female', 4, 90, 'Infinity', 5),
    ('立定跳远(厘米)','male', 4.5, 40, 54, 1),
    ('立定跳远(厘米)','male', 4.5, 55, 72, 2),
    ('立定跳远(厘米)','male', 4.5, 73, 89, 3),
    ('立定跳远(厘米)','male', 4.5, 90, 102, 4),
    ('立定跳远(厘米)','male', 4.5, 103, 'Infinity', 5),
    ('立定跳远(厘米)','female', 4.5, 40, 49, 1),
    ('立定跳远(厘米)','female', 4.5, 50, 67, 2),
    ('立定跳远(厘米)','female', 4.5, 68, 80, 3),
    ('立定跳远(厘米)','female', 4.5, 81, 96, 4),
    ('立定跳远(厘米)','female', 4.5, 97, 'Infinity', 5),
    ('立定跳远(厘米)','male', 5, 50, 64, 1),
    ('立定跳远(厘米)','male', 5, 65, 79, 2),
    ('立定跳远(厘米)','male', 5, 80, 95, 3),
    ('立定跳远(厘米)','male', 5, 96, 110, 4),
    ('立定跳远(厘米)','male', 5, 111, 'Infinity', 5),
    ('立定跳远(厘米)','female', 5, 50, 59, 1),
    ('立定跳远(厘米)','female', 5, 60, 74, 2),
    ('立定跳远(厘米)','female', 5, 75, 88, 3),
    ('立定跳远(厘米)','female', 5, 89, 102, 4),
    ('立定跳远(厘米)','female', 5, 103, 'Infinity', 5),
    ('立定跳远(厘米)','male', 5.5, 56, 69, 1),
    ('立定跳远(厘米)','male', 5.5, 70, 89, 2),
    ('立定跳远(厘米)','male', 5.5, 90, 102, 3),
    ('立定跳远(厘米)','male', 5.5, 103, 119, 4),
    ('立定跳远(厘米)','male', 5.5, 120, 'Infinity', 5),
    ('立定跳远(厘米)','female', 5.5, 54, 65, 1),
    ('立定跳远(厘米)','female', 5.5, 66, 81, 2),
    ('立定跳远(厘米)','female', 5.5, 82, 95, 3),
    ('立定跳远(厘米)','female', 5.5, 96, 109, 4),
    ('立定跳远(厘米)','female', 5.5, 110, 'Infinity', 5),
    ('立定跳远(厘米)','male', 6, 61, 78, 1),
    ('立定跳远(厘米)','male', 6, 79, 94, 2),
    ('立定跳远(厘米)','male', 6, 95, 110, 3),
    ('立定跳远(厘米)','male', 6, 111, 127, 4),
    ('立定跳远(厘米)','male', 6, 128, 'Infinity', 5),
    ('立定跳远(厘米)','female', 6, 60, 70, 1),
    ('立定跳远(厘米)','female', 6, 71, 86, 2),
    ('立定跳远(厘米)','female', 6, 87, 100, 3),
    ('立定跳远(厘米)','female', 6, 101, 116, 4),
    ('立定跳远(厘米)','female', 6, 117, 'Infinity', 5);

INSERT INTO "standard_scale_score"("name","gender","age","min","max","score") VALUES
    ('网球掷远(米)','male', 3, 1.5, 1.5, 1),
    ('网球掷远(米)','male', 3, 2.0, 2.5, 2),
    ('网球掷远(米)','male', 3, 3.0, 3.5, 3),
    ('网球掷远(米)','male', 3, 4.0, 5.5, 4),
    ('网球掷远(米)','male', 3, 6.0, 'Infinity', 5),
    ('网球掷远(米)','female', 3, 1, 1, 1),
    ('网球掷远(米)','female', 3, 1.5, 2, 2),
    ('网球掷远(米)','female', 3, 2.5, 3, 3),
    ('网球掷远(米)','female', 3, 3.5, 5, 4),
    ('网球掷远(米)','female', 3, 5.5, 'Infinity', 5),
    ('网球掷远(米)','male', 3.5, 1.5, 1.5, 1),
    ('网球掷远(米)','male', 3.5, 2.0, 2.5, 2),
    ('网球掷远(米)','male', 3.5, 3.0, 4.0, 3),
    ('网球掷远(米)','male', 3.5, 4.5, 5.5, 4),
    ('网球掷远(米)','male', 3.5, 6.0, 'Infinity', 5),
    ('网球掷远(米)','female', 3.5, 1.5, 1.5, 1),
    ('网球掷远(米)','female', 3.5, 2.0, 2.5, 2),
    ('网球掷远(米)','female', 3.5, 3.0, 3.5, 3),
    ('网球掷远(米)','female', 3.5, 40, 5.0, 4),
    ('网球掷远(米)','female', 3.5, 5.5, 'Infinity', 5),
    ('网球掷远(米)','male', 4, 2.0, 2.5, 1),
    ('网球掷远(米)','male', 4, 3.0, 3.5, 2),
    ('网球掷远(米)','male', 4, 4.0, 4.5, 3),
    ('网球掷远(米)','male', 4, 5.0, 6.0, 4),
    ('网球掷远(米)','male', 4, 6.5, 'Infinity', 5),
    ('网球掷远(米)','female', 4, 2.0, 2.0, 1),
    ('网球掷远(米)','female', 4, 2.5, 3, 2),
    ('网球掷远(米)','female', 4, 3.5, 4.0, 3),
    ('网球掷远(米)','female', 4, 4.5, 5.0, 4),
    ('网球掷远(米)','female', 4, 5.5, 'Infinity', 5),
    ('网球掷远(米)','male', 4.5, 2.5, 2.5, 1),
    ('网球掷远(米)','male', 4.5, 3.0, 4.0, 2),
    ('网球掷远(米)','male', 4.5, 4.5, 6.0, 3),
    ('网球掷远(米)','male', 4.5, 6.5, 8.0, 4),
    ('网球掷远(米)','male', 4.5, 8.5, 'Infinity', 5),
    ('网球掷远(米)','female', 4.5, 2, 2, 1),
    ('网球掷远(米)','female', 4.5, 2.5, 3.0, 2),
    ('网球掷远(米)','female', 4.5, 3.5, 4.0, 3),
    ('网球掷远(米)','female', 4.5, 4.5, 5.5, 4),
    ('网球掷远(米)','female', 4.5, 6.0, 'Infinity', 5),
    ('网球掷远(米)','male', 5.0, 3, 3.5, 1),
    ('网球掷远(米)','male', 5.0, 4.0, 5.0, 2),
    ('网球掷远(米)','male', 5.0, 5.5, 7.0, 3),
    ('网球掷远(米)','male', 5.0, 7.5, 9.0, 4),
    ('网球掷远(米)','male', 5.0, 9.5, 'Infinity', 5),
    ('网球掷远(米)','female', 5.0, 2.5, 3, 1),
    ('网球掷远(米)','female', 5.0, 3.5, 4.0, 2),
    ('网球掷远(米)','female', 5.0, 4.5, 5.5, 3),
    ('网球掷远(米)','female', 5.0, 6.0, 8.5, 4),
    ('网球掷远(米)','female', 5.0, 9, 'Infinity', 5),
    ('网球掷远(米)','male', 5.5, 3, 3.5, 1),
    ('网球掷远(米)','male', 5.5, 4.0, 5.5, 2),
    ('网球掷远(米)','male', 5.5, 6.0, 7.5, 3),
    ('网球掷远(米)','male', 5.5, 8.0, 10.0, 4),
    ('网球掷远(米)','male', 5.5, 10.5, 'Infinity', 5),
    ('网球掷远(米)','female', 5.5, 3, 3, 1),
    ('网球掷远(米)','female', 5.5, 3.5, 4.5, 2),
    ('网球掷远(米)','female', 5.5, 5.0, 6.0, 3),
    ('网球掷远(米)','female', 5.5, 6.5, 8.5, 4),
    ('网球掷远(米)','female', 5.5, 9, 'Infinity', 5),
    ('网球掷远(米)','male', 6, 3.5, 4, 1),
    ('网球掷远(米)','male', 6, 4.5, 6.5, 2),
    ('网球掷远(米)','male', 6, 7.0, 9.0, 3),
    ('网球掷远(米)','male', 6, 9.5, 12.0, 4),
    ('网球掷远(米)','male', 6, 12.5, 'Infinity', 5),
    ('网球掷远(米)','female', 6, 3, 3, 1),
    ('网球掷远(米)','female', 6, 3.5, 4.5, 2),
    ('网球掷远(米)','female', 6, 5.0, 6.0, 3),
    ('网球掷远(米)','female', 6, 6.5, 8.0, 4),
    ('网球掷远(米)','female', 6, 9, 'Infinity', 5);


INSERT INTO "standard_scale_score"("name","gender","age","min","max","score") VALUES
    ('双脚连续跳(秒)','male', 3, 19.7, 25, 1),
    ('双脚连续跳(秒)','male', 3, 13.1, 19.6, 2),
    ('双脚连续跳(秒)','male', 3, 9.2, 13.0, 3),
    ('双脚连续跳(秒)','male', 3, 6.6, 9.1, 4),
    ('双脚连续跳(秒)','male', 3, 0, 6.5, 5),
    ('双脚连续跳(秒)','female', 3, 20.1, 25.9, 1),
    ('双脚连续跳(秒)','female', 3, 13.5, 20, 2),
    ('双脚连续跳(秒)','female', 3, 9.8, 13.4, 3),
    ('双脚连续跳(秒)','female', 3, 7.1, 9.7, 4),
    ('双脚连续跳(秒)','female', 3, 0, 7.0, 5),
    ('双脚连续跳(秒)','male', 3.5, 17, 21.8, 1),
    ('双脚连续跳(秒)','male', 3.5, 11.2, 16.9, 2),
    ('双脚连续跳(秒)','male', 3.5, 8.3, 11.1, 3),
    ('双脚连续跳(秒)','male', 3.5, 6.1, 8.2, 4),
    ('双脚连续跳(秒)','male', 3.5, 0, 6, 5),
    ('双脚连续跳(秒)','female', 3.5, 17.1, 21.9, 1),
    ('双脚连续跳(秒)','female', 3.5, 11.3, 17, 2),
    ('双脚连续跳(秒)','female', 3.5, 8.5, 11.2, 3),
    ('双脚连续跳(秒)','female', 3.5, 6.2, 8.4, 4),
    ('双脚连续跳(秒)','female', 3.5, 0, 6.1, 5),
    ('双脚连续跳(秒)','male', 4, 13.2, 17, 1),
    ('双脚连续跳(秒)','male', 4, 9.2, 13.1, 2),
    ('双脚连续跳(秒)','male', 4, 7.1, 9.1, 3),
    ('双脚连续跳(秒)','male', 4, 5.6, 7, 4),
    ('双脚连续跳(秒)','male', 4, 0, 5.5, 5),
    ('双脚连续跳(秒)','female', 4, 13.5, 17.2, 1),
    ('双脚连续跳(秒)','female', 4, 9.6, 13.4, 2),
    ('双脚连续跳(秒)','female', 4, 7.4, 9.5, 3),
    ('双脚连续跳(秒)','female', 4, 5.9, 7.3, 4),
    ('双脚连续跳(秒)','female', 4, 0, 5.8, 5),
    ('双脚连续跳(秒)','male', 4.5, 11.3, 14.5, 1),
    ('双脚连续跳(秒)','male', 4.5, 8.2, 11.2, 2),
    ('双脚连续跳(秒)','male', 4.5, 6.5, 8.1, 3),
    ('双脚连续跳(秒)','male', 4.5, 5.3, 6.4, 4),
    ('双脚连续跳(秒)','male', 4.5, 0, 5.2, 5),
    ('双脚连续跳(秒)','female', 4.5, 12, 14.9, 1),
    ('双脚连续跳(秒)','female', 4.5, 8.6, 11.9, 2),
    ('双脚连续跳(秒)','female', 4.5, 6.8, 8.5, 3),
    ('双脚连续跳(秒)','female', 4.5, 5.5, 6.7, 4),
    ('双脚连续跳(秒)','female', 4.5, 0, 5.4, 5),
    ('双脚连续跳(秒)','male', 5, 9.9, 12.5, 1),
    ('双脚连续跳(秒)','male', 5, 7.3, 9.8, 2),
    ('双脚连续跳(秒)','male', 5, 6.0, 7.2, 3),
    ('双脚连续跳(秒)','male', 5, 5.1, 5.9, 4),
    ('双脚连续跳(秒)','male', 5, 0, 5.0, 5),
    ('双脚连续跳(秒)','female', 5, 10.1, 12.7, 1),
    ('双脚连续跳(秒)','female', 5, 7.6, 10, 2),
    ('双脚连续跳(秒)','female', 5, 6.2, 7.5, 3),
    ('双脚连续跳(秒)','female', 5, 5.2, 6.1, 4),
    ('双脚连续跳(秒)','female', 5, 0, 5.1, 5),
    ('双脚连续跳(秒)','male', 5.5, 9.4, 11.9, 1),
    ('双脚连续跳(秒)','male', 5.5, 6.9, 9.3, 2),
    ('双脚连续跳(秒)','male', 5.5, 5.7, 6.8, 3),
    ('双脚连续跳(秒)','male', 5.5, 4.9, 5.6, 4),
    ('双脚连续跳(秒)','male', 5.5, 0, 4.8, 5),
    ('双脚连续跳(秒)','female', 5.5, 9.3, 11.5, 1),
    ('双脚连续跳(秒)','female', 5.5, 7, 9.2, 2),
    ('双脚连续跳(秒)','female', 5.5, 5.8, 6.9, 3),
    ('双脚连续跳(秒)','female', 5.5, 4.9, 5.7, 4),
    ('双脚连续跳(秒)','female', 5.5, 0, 4.8, 5),
    ('双脚连续跳(秒)','male', 6, 8.3, 10.4, 1),
    ('双脚连续跳(秒)','male', 6, 6.2, 8.2, 2),
    ('双脚连续跳(秒)','male', 6, 5.2, 6.1, 3),
    ('双脚连续跳(秒)','male', 6, 4.4, 5.1, 4),
    ('双脚连续跳(秒)','male', 6, 0, 4.3, 5),
    ('双脚连续跳(秒)','female', 6, 8.4, 10.5, 1),
    ('双脚连续跳(秒)','female', 6, 6.3, 8.3, 2),
    ('双脚连续跳(秒)','female', 6, 5.3, 6.2, 3),
    ('双脚连续跳(秒)','female', 6, 4.6, 5.2, 4),
    ('双脚连续跳(秒)','female', 6, 0, 4.5, 5);

INSERT INTO "standard_scale_score"("name","gender","age","min","max","score") VALUES
    ('坐位体前屈(厘米)','male', 3, 2.9, 4.8, 1),
    ('坐位体前屈(厘米)','male', 3, 4.9, 8.5, 2),
    ('坐位体前屈(厘米)','male', 3, 8.6, 11.6, 3),
    ('坐位体前屈(厘米)','male', 3, 11.7, 14.9, 4),
    ('坐位体前屈(厘米)','male', 3, 15, 'Infinity', 5),
    ('坐位体前屈(厘米)','female', 3, 2.7, 4.6, 1),
    ('坐位体前屈(厘米)','female', 3, 4.7, 9.9, 2),
    ('坐位体前屈(厘米)','female', 3, 10, 12.9, 3),
    ('坐位体前屈(厘米)','female', 3, 13, 15.9, 4),
    ('坐位体前屈(厘米)','female', 3, 16, 'Infinity', 5),
    ('坐位体前屈(厘米)','male', 3.5, 2.7, 4.6, 1),
    ('坐位体前屈(厘米)','male', 3.5, 4.7, 8.7, 2),
    ('坐位体前屈(厘米)','male', 3.5, 8.8, 11.5, 3),
    ('坐位体前屈(厘米)','male', 3.5, 11.6, 14.9, 4),
    ('坐位体前屈(厘米)','male', 3.5, 15, 'Infinity', 5),
    ('坐位体前屈(厘米)','female', 3.5, 3.5, 6.2, 1),
    ('坐位体前屈(厘米)','female', 3.5, 6.3, 9.9, 2),
    ('坐位体前屈(厘米)','female', 3.5, 10.0, 12.9, 3),
    ('坐位体前屈(厘米)','female', 3.5, 13, 15.9, 4),
    ('坐位体前屈(厘米)','female', 3.5, 16, 'Infinity', 5),
    ('坐位体前屈(厘米)','male', 4, 2.4, 4.4, 1),
    ('坐位体前屈(厘米)','male', 4, 4.5, 8.5, 2),
    ('坐位体前屈(厘米)','male', 4, 8.6, 11.4, 3),
    ('坐位体前屈(厘米)','male', 4, 11.5, 14.9, 4),
    ('坐位体前屈(厘米)','male', 4, 15, 'Infinity', 5),
    ('坐位体前屈(厘米)','female', 4, 3, 5.9, 1),
    ('坐位体前屈(厘米)','female', 4, 6, 9.9, 2),
    ('坐位体前屈(厘米)','female', 4, 10.0, 12.9, 3),
    ('坐位体前屈(厘米)','female', 4, 13, 15.9, 4),
    ('坐位体前屈(厘米)','female', 4, 16, 'Infinity', 5),
    ('坐位体前屈(厘米)','male', 4.5, 1.8, 4.1, 1),
    ('坐位体前屈(厘米)','male', 4.5, 4.2, 7.9, 2),
    ('坐位体前屈(厘米)','male', 4.5, 8, 10.9, 3),
    ('坐位体前屈(厘米)','male', 4.5, 11, 14.4, 4),
    ('坐位体前屈(厘米)','male', 4.5, 14.5, 'Infinity', 5),
    ('坐位体前屈(厘米)','female', 4.5, 3, 5.9, 1),
    ('坐位体前屈(厘米)','female', 4.5, 6, 9.9, 2),
    ('坐位体前屈(厘米)','female', 4.5, 10.0, 12.9, 3),
    ('坐位体前屈(厘米)','female', 4.5, 13, 15.9, 4),
    ('坐位体前屈(厘米)','female', 4.5, 16, 'Infinity', 5),
    ('坐位体前屈(厘米)','male', 5, 1.1, 3.4, 1),
    ('坐位体前屈(厘米)','male', 5, 3.5, 7.5, 2),
    ('坐位体前屈(厘米)','male', 5, 7.6, 10.9, 3),
    ('坐位体前屈(厘米)','male', 5, 11, 14.4, 4),
    ('坐位体前屈(厘米)','male', 5, 14.5, 'Infinity', 5),
    ('坐位体前屈(厘米)','female', 5, 3, 5.4, 1),
    ('坐位体前屈(厘米)','female', 5, 5.5, 9.6, 2),
    ('坐位体前屈(厘米)','female', 5, 9.7, 13.1, 3),
    ('坐位体前屈(厘米)','female', 5, 13.2, 16.6, 4),
    ('坐位体前屈(厘米)','female', 5, 16.7, 'Infinity', 5),
    ('坐位体前屈(厘米)','male', 5.5, 1, 3.2, 1),
    ('坐位体前屈(厘米)','male', 5.5, 3.3, 7.5, 2),
    ('坐位体前屈(厘米)','male', 5.5, 7.6, 10.9, 3),
    ('坐位体前屈(厘米)','male', 5.5, 11, 14.4, 4),
    ('坐位体前屈(厘米)','male', 5.5, 14.5, 'Infinity', 5),
    ('坐位体前屈(厘米)','female', 5.5, 3, 5.4, 1),
    ('坐位体前屈(厘米)','female', 5.5, 5.5, 9.6, 2),
    ('坐位体前屈(厘米)','female', 5.5, 9.7, 12.9, 3),
    ('坐位体前屈(厘米)','female', 5.5, 13, 16.7, 4),
    ('坐位体前屈(厘米)','female', 5.5, 16.8, 'Infinity', 5),
    ('坐位体前屈(厘米)','male', 6, 1, 3.1, 1),
    ('坐位体前屈(厘米)','male', 6, 3.2, 7, 2),
    ('坐位体前屈(厘米)','male', 6, 7.1, 10.4, 3),
    ('坐位体前屈(厘米)','male', 6, 10.5, 14.4, 4),
    ('坐位体前屈(厘米)','male', 6, 14.5, 'Infinity', 5),
    ('坐位体前屈(厘米)','female', 6, 3, 5.4, 1),
    ('坐位体前屈(厘米)','female', 6, 5.5, 9.5, 2),
    ('坐位体前屈(厘米)','female', 6, 9.6, 12.9, 3),
    ('坐位体前屈(厘米)','female', 6, 13, 16.7, 4),
    ('坐位体前屈(厘米)','female', 6, 16.8, 'Infinity', 5);


INSERT INTO "standard_scale_score"("name","gender","age","min","max","score") VALUES
    ('走平衡木(秒)','male', 3, 30.1, 48.5, 1),
    ('走平衡木(秒)','male', 3, 16.9, 30, 2),
    ('走平衡木(秒)','male', 3, 10.6, 16.8, 3),
    ('走平衡木(秒)','male', 3, 6.6, 10.5, 4),
    ('走平衡木(秒)','male', 3, 0, 6.5, 5),
    ('走平衡木(秒)','female', 3, 32.5, 49.8, 1),
    ('走平衡木(秒)','female', 3, 17.4, 32.4, 2),
    ('走平衡木(秒)','female', 3, 10.8, 17.3, 3),
    ('走平衡木(秒)','female', 3, 6.9, 10.7, 4),
    ('走平衡木(秒)','female', 3, 0, 6.8, 5),
    ('走平衡木(秒)','male', 3.5, 27.1, 41.1, 1),
    ('走平衡木(秒)','male', 3.5, 15.1, 27, 2),
    ('走平衡木(秒)','male', 3.5, 9.4, 15, 3),
    ('走平衡木(秒)','male', 3.5, 5.9, 9.3, 4),
    ('走平衡木(秒)','male', 3.5, 0, 5.8, 5),
    ('走平衡木(秒)','female', 3.5, 27.5, 40.4, 1),
    ('走平衡木(秒)','female', 3.5, 15.1, 27.4, 2),
    ('走平衡木(秒)','female', 3.5, 9.7, 15, 3),
    ('走平衡木(秒)','female', 3.5, 6.1, 9.6, 4),
    ('走平衡木(秒)','female', 3.5, 0, 6, 5),
    ('走平衡木(秒)','male', 4, 21.6, 33.2, 1),
    ('走平衡木(秒)','male', 4, 11.6, 21.5, 2),
    ('走平衡木(秒)','male', 4, 7.4, 11.5, 3),
    ('走平衡木(秒)','male', 4, 4.9, 7.3, 4),
    ('走平衡木(秒)','male', 4, 0, 4.8, 5),
    ('走平衡木(秒)','female', 4, 22.6, 32.2, 1),
    ('走平衡木(秒)','female', 4, 12.3, 22.5, 2),
    ('走平衡木(秒)','female', 4, 8.2, 12.2, 3),
    ('走平衡木(秒)','female', 4, 5.8, 8.1, 4),
    ('走平衡木(秒)','female', 4, 0, 5.7, 5),
    ('走平衡木(秒)','male', 4.5, 17.9, 28.4, 1),
    ('走平衡木(秒)','male', 4.5, 9.7, 17.8, 2),
    ('走平衡木(秒)','male', 4.5, 6.3, 9.6, 3),
    ('走平衡木(秒)','male', 4.5, 4.3, 6.2, 4),
    ('走平衡木(秒)','male', 4.5, 0, 4.2, 5),
    ('走平衡木(秒)','female', 4.5, 18.5, 26.5, 1),
    ('走平衡木(秒)','female', 4.5, 10.2, 18.4, 2),
    ('走平衡木(秒)','female', 4.5, 7.0, 10.1, 3),
    ('走平衡木(秒)','female', 4.5, 4.7, 6.9, 4),
    ('走平衡木(秒)','female', 4.5, 0, 4.6, 5),
    ('走平衡木(秒)','male', 5, 14.1, 22.2, 1),
    ('走平衡木(秒)','male', 5, 7.9, 14, 2),
    ('走平衡木(秒)','male', 5, 5.3, 7.8, 3),
    ('走平衡木(秒)','male', 5, 3.7, 5.2, 4),
    ('走平衡木(秒)','male', 5, 0, 3.6, 5),
    ('走平衡木(秒)','female', 5, 14.1, 23.7, 1),
    ('走平衡木(秒)','female', 5, 8.3, 14, 2),
    ('走平衡木(秒)','female', 5, 5.8, 8.2, 3),
    ('走平衡木(秒)','female', 5, 4.1, 5.7, 4),
    ('走平衡木(秒)','female', 5, 0, 4, 5),
    ('走平衡木(秒)','male', 5.5, 12.1, 19.2, 1),
    ('走平衡木(秒)','male', 5.5, 6.8, 12, 2),
    ('走平衡木(秒)','male', 5.5, 4.6, 6.7, 3),
    ('走平衡木(秒)','male', 5.5, 3.3, 4.5, 4),
    ('走平衡木(秒)','male', 5.5, 0, 3.2, 5),
    ('走平衡木(秒)','female', 5.5, 12.6, 20.1, 1),
    ('走平衡木(秒)','female', 5.5, 7.5, 12.5, 2),
    ('走平衡木(秒)','female', 5.5, 5.1, 7.4, 3),
    ('走平衡木(秒)','female', 5.5, 3.6, 5.0, 4),
    ('走平衡木(秒)','female', 5.5, 0, 3.5, 5),
    ('走平衡木(秒)','male', 6, 9.4, 16, 1),
    ('走平衡木(秒)','male', 6, 5.4, 9.3, 2),
    ('走平衡木(秒)','male', 6, 3.8, 5.3, 3),
    ('走平衡木(秒)','male', 6, 2.7, 3.7, 4),
    ('走平衡木(秒)','male', 6, 0, 2.6, 5),
    ('走平衡木(秒)','female', 6, 10.8, 17, 1),
    ('走平衡木(秒)','female', 6, 6.2, 10.7, 2),
    ('走平衡木(秒)','female', 6, 4.1, 6.1, 3),
    ('走平衡木(秒)','female', 6, 3.0, 4.2, 4),
    ('走平衡木(秒)','female', 6, 0, 2.9, 5);


CREATE TABLE "standard_scale_hw_score"(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "gender" VARCHAR(32) NOT NULL, -- 性别
    "height_min" FLOAT NOT NULL,
    "height_max" FLOAT NOT NULL,
    "weight_min" FLOAT NOT NULL,
    "weight_max" FLOAT NOT NULL,
    "score" FLOAT NOT NULL, -- 分数
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE("gender","height_min","height_max","weight_min","weight_max")
);
CREATE INDEX ON "standard_scale_hw_score"("gender");
CREATE INDEX ON "standard_scale_hw_score"("height_min");
CREATE INDEX ON "standard_scale_hw_score"("height_max");
CREATE INDEX ON "standard_scale_hw_score"("weight_min");
CREATE INDEX ON "standard_scale_hw_score"("weight_max");
