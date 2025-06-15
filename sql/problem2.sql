-- 요구사항:
-- 1. ETD 기간: 2025년 5월 1일 ~ 2025년 5월 11일
-- 2. POL(출발항) + 컨테이너 타입(CNTR_TYPE) 별로 집계
-- 3. 컨테이너 수량(CNTR 개수) 및 총 중량 합계(CNTR_WGT) 조회
-- 4. 결과 컬럼: POL_CD, CNTR_TYPE, CNTR_COUNT, TOTAL_WGT
-- 5. 정렬: POL_CD, CNTR_TYPE

-- 문제 2 정답
SELECT
    m.POL_CD,
    c.CNTR_TYPE,
    COUNT(c.CNTR_NO) AS CNTR_COUNT,
    SUM(c.CNTR_WGT) AS TOTAL_WGT
FROM
    FMS_HBL_MST m
JOIN
    FMS_HBL_CNTR c ON m.HBL_NO = c.HBL_NO
WHERE
    m.ETD BETWEEN '20250501' AND '20250511'
    AND c.CNTR_NO IS NOT NULL
GROUP BY
    m.POL_CD, c.CNTR_TYPE
ORDER BY
    m.POL_CD, c.CNTR_TYPE;
