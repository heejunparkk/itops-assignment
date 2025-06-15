-- 요구사항:
-- 1. HOUSE B/L (HBL_NO) 를 기준으로 각 B/L 별 컨테이너 수량(CNTR 개수) 을 집계
-- 2. 가장 많은 수량을 가진 B/L 1건을 조회
-- 3. 컨테이너 수량은 CNTR_NO 기준으로 COUNT
-- 4. 동일 수량 시 ETD 빠른 순으로 우선 선택
-- 5. 결과 컬럼: HBL_NO, CNTR_COUNT, ETD
-- 6. 정렬: CNTR_COUNT DESC, ETD ASC

-- 문제 1 정답
WITH CONTAINER_COUNTS AS (
    SELECT 
        h.HBL_NO,
        COUNT(DISTINCT c.CNTR_NO) AS CNTR_COUNT,
        h.ETD
    FROM 
        FMS_HBL_MST h
    JOIN 
        FMS_HBL_CNTR c ON h.HBL_NO = c.HBL_NO
    WHERE 
        c.CNTR_NO IS NOT NULL
    GROUP BY 
        h.HBL_NO, h.ETD
),
RANKED_RESULTS AS (
    SELECT 
        HBL_NO,
        CNTR_COUNT,
        ETD,
        RANK() OVER (ORDER BY CNTR_COUNT DESC, ETD ASC) AS RANK_NUM
    FROM 
        CONTAINER_COUNTS
)
SELECT 
    HBL_NO,
    CNTR_COUNT,
    ETD
FROM 
    RANKED_RESULTS
WHERE 
    RANK_NUM = 1;
