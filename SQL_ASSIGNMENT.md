# ItOps 개발자 채용 과제 - Oracle SQL 쿼리 문제

## 개요

이 문제는 물류 시스템의 HOUSE B/L과 컨테이너 데이터를 기반으로 한 Oracle SQL 쿼리 작성 문제입니다.

## 제출 방법

- 각 문제의 답안을 별도의 SQL 파일로 작성해 주세요
- 파일명: `problem1.sql`, `problem2.sql`
- GitHub 레포지토리의 `sql/` 폴더에 업로드

## 공통 KEY: HOUSE B/L (HBL_NO)

두 테이블은 모두 HBL_NO 를 기준으로 연계 (JOIN, GROUP BY 등에서 사용)

## 테이블 설명

### 1. HOUSE B/L MASTER 테이블 (FMS_HBL_MST)

| Column name | Type         | Nullable | 설명                     |
| ----------- | ------------ | -------- | ------------------------ |
| HBL_NO      | varchar2(20) |          | HOUSE B/L 번호 (PK 역할) |
| POL_CD      | varchar2(5)  | Y        | 출발항 코드              |
| POL_NM      | varchar2(50) | Y        | 출발항 명                |
| POD_CD      | varchar2(5)  | Y        | 도착항 코드              |
| POD_NM      | varchar2(50) | Y        | 도착항 명                |
| ONBD_YMD    | varchar2(8)  | Y        | On board 일자            |
| VSL_CD      | varchar2(20) | Y        | 선명 코드                |
| VSL         | varchar2(35) | Y        | 선명                     |
| VOY         | varchar2(10) | Y        | 항차                     |
| ETD         | varchar2(8)  | Y        | 출항일자                 |
| ETA         | varchar2(8)  | Y        | 입항일자                 |

### 2. 컨테이너 테이블 (FMS_HBL_CNTR)

| Column name | Type         | Nullable | 설명                     |
| ----------- | ------------ | -------- | ------------------------ |
| HBL_NO      | varchar2(20) |          | HOUSE B/L 번호 (FK 역할) |
| CNTR_SEQ    | number(4)    |          | 컨테이너 순번            |
| CNTR_NO     | varchar2(11) | Y        | 컨테이너 번호            |
| CNTR_TYPE   | varchar2(4)  | Y        | 컨테이너 타입            |
| CNTR_WGT    | number(15,3) |          | 컨테이너 중량            |

## 예시 데이터

### FMS_HBL_MST 테이블

| HBL_NO           | POL_CD | POL_NM       | POD_CD | POD_NM       | ETD      | ETA      |
| ---------------- | ------ | ------------ | ------ | ------------ | -------- | -------- |
| PSECSOSA25050009 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250511 | 20250512 |
| PSECSOSA25050001 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250503 | 20250505 |
| PSECSOSA25050002 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250508 | 20250509 |
| PSECSOSA25050005 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250512 | 20250514 |
| PSECSOSA25050055 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250513 | 20250514 |
| PSECSOSA25050061 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250513 | 20250514 |
| PSECSOSA25050021 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250511 | 20250512 |
| PSECSOSA25050025 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250508 | 20250509 |
| PSECSOSA2505545D | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250501 | 20250502 |
| PSECSOSA25050003 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250513 | 20250514 |
| PSECSOSA25050024 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250506 | 20250507 |
| PSECSOSA25050036 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250503 | 20250505 |
| PSECSOSA25050089 | KRPUS  | BUSAN, KOREA | JPOSA  | OSAKA, JAPAN | 20250513 | 20250514 |

### FMS_HBL_CNTR 테이블

| HBL_NO           | CNTR_SEQ | CNTR_NO     | CNTR_TYPE | SEAL_NO1  | PKG     | PKG_UNIT_CD | CNTR_WGT | CNTR_MSRMT |
| ---------------- | -------- | ----------- | --------- | --------- | ------- | ----------- | -------- | ---------- |
| PSECSOSA25050001 | 1        | HLHU8444991 | 45GP      | 080277    | 454.00  | CT          | 3793.00  | 57.85      |
| PSECSOSA25050002 | 1        | KDCU5214334 | 45GP      | 587689    | 533.00  | CT          | 4159.00  | 67.92      |
| PSECSOSA25050003 | 1        | PHRU9131666 | 45GP      | 621368    | 533.00  | CT          | 4159.00  | 67.92      |
| PSECSOSA25050005 | 1        | SKHU6314194 | 45GP      | 080217    | 533.00  | CT          | 4159.00  | 67.92      |
| PSECSOSA25050009 | 1        | PHRU8800242 | 45GP      | PST531456 | 800.00  | CT          | 4320.00  | 28.00      |
| PSECSOSA25050021 | 1        | PSEU2110359 | 22GP      | 531799    | 3.00    | PK          | 2440.00  | 13.86      |
| PSECSOSA25050024 | 1        | PHRU5623380 | 22RE      | 639348    | 780.00  | PK          | 8381.00  | 24.89      |
| PSECSOSA25050025 | 1        | PHRU5619061 | 22RE      | 605228    | 961.00  | PK          | 11187.30 | 24.96      |
| PSECSOSA25050036 | 1        | STXU4555540 | 45GP      | 090736    | 54.00   | PL          | 13446.00 | 50.00      |
| PSECSOSA25050036 | 2        | STXU4536637 | 45GP      | 090737    | 54.00   | PL          | 13446.00 | 50.00      |
| PSECSOSA25050055 | 1        | PHRU2207836 | 22GP      | 639153    | 10.00   | PL          | 9408.00  | 25.77      |
| PSECSOSA25050061 | 1        | PHRU5621222 | 22RE      | 531757    | 830.00  | PK          | 9569.00  | 24.88      |
| PSECSOSA25050089 | 1        | PSEU5101556 | 45GP      | 621381    | 16.00   | PL          | 19896.26 | 48.16      |
| PSECSOSA25050089 | 2        | PHRU8600054 | 45GP      | 507244    | 16.00   | PL          | 20018.56 | 48.16      |
| PSECSOSA25050089 | 3        | PHRU8601220 | 45GP      | 507815    | 16.00   | PL          | 20018.56 | 48.16      |
| PSECSOSA2505545D | 1        | PHRU5210200 | 45RE      | PST531430 | 5501.00 | CT          | 13986.92 | 32.49      |

## 문제 1

### 컨테이너 수량이 가장 많은 B/L 구하기

- HOUSE B/L (HBL_NO) 를 기준으로
- 각 B/L 별 컨테이너 수량(CNTR 개수) 을 집계하고
- 가장 많은 수량을 가진 B/L 1건을 조회하는 SQL문을 작성하세요.

#### 조건

- 컨테이너 수량은 CNTR_NO 기준으로 COUNT
- 동일 수량 시 ETD 빠른 순으로 우선 선택

#### 결과 컬럼:

- HBL_NO
- CNTR_COUNT
- ETD

#### 정렬 조건

- CNTR_COUNT가 많고 ETD 오름차순으로 정렬

#### 출력 예시 (예상 결과 예시)

| HBL_NO           | CNTR_COUNT | ETD      |
| ---------------- | ---------- | -------- |
| PSECSOSA25050009 | 3          | 20250511 |

## 문제 2

### ETD 기간 내 POL별 컨테이너 타입별 컨테이너 수량 및 중량 합계 조회하기 (ETD는 필터만 적용)

- 지정한 ETD 기간 (2025년 5월 1일 ~ 2025년 5월 11일) 내에서
- POL(출발항) + 컨테이너 타입(CNTR_TYPE) 별로
- 컨테이너 수량(CNTR 개수) 및 총 중량 합계(CNTR_WGT) 를 조회하는 SQL문을 작성하세요.

#### 결과 컬럼

| POL_CD | CNTR_TYPE | CNTR_COUNT | TOTAL_WGT |

#### 정렬 조건

ORDER BY POL_CD, CNTR_TYPE

#### 출력 예시 (예상 결과 예시)

| POL_CD | CNTR_TYPE | CNTR_COUNT | TOTAL_WGT |
| ------ | --------- | ---------- | --------- |
| KRPUS  | 22GP      | 1          | 2440.00   |
| KRPUS  | 22RE      | 2          | 19568.30  |
| KRPUS  | 45GP      | 5          | 39164.00  |
| KRPUS  | 45RE      | 1          | 13986.92  |
