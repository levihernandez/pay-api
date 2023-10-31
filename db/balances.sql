WITH Transfers AS (
    SELECT
        L.id,
        L.sender_uuid,
        S.balance AS sender_before,
        S.balance - L.amount AS sender_after,
        L.receiver_uuid,
        R.balance AS receiver_before,
        R.balance + L.amount AS receiver_after,
        L.amount,
        L.created_at
    FROM
        ledger L
        JOIN account_types S ON L.sender_uuid = S.customer_uuid
        JOIN account_types R ON L.receiver_uuid = R.customer_uuid
    WHERE
        S.account_type = L.account_type AND R.account_type = L.account_type
		-- AND L.sender_uuid = '8b6ca348-5525-4a16-8da1-e031e5790909'
    ORDER BY
        L.created_at
)
SELECT
    T.id,
    C1.name AS sender_name,
    T.sender_before,
    T.sender_after,
    T.amount,
    C2.name AS receiver_name,
    T.receiver_before,
    T.receiver_after,
    T.created_at
FROM
    Transfers T
    JOIN customers C1 ON T.sender_uuid = C1.uuid
    JOIN customers C2 ON T.receiver_uuid = C2.uuid
	WHERE T.id = 47
ORDER BY
    T.created_at;

/*
id|sender_name|sender_before|sender_after|amount|receiver_name|receiver_before|receiver_after|created_at         |
--+-----------+-------------+------------+------+-------------+---------------+--------------+-------------------+
47|Robert     |      10074.0|     10019.0|  55.0|Brian        |          709.0|         764.0|2023-10-31 01:14:35|
47|Robert     |      10074.0|     10019.0|  55.0|Brian        |         4973.0|        5028.0|2023-10-31 01:14:35|
47|Robert     |      10074.0|     10019.0|  55.0|Brian        |         7268.0|        7323.0|2023-10-31 01:14:35|
47|Robert     |      10074.0|     10019.0|  55.0|Brian        |         8926.0|        8981.0|2023-10-31 01:14:35|
47|Robert     |       1589.0|      1534.0|  55.0|Brian        |          709.0|         764.0|2023-10-31 01:14:35|
47|Robert     |       1589.0|      1534.0|  55.0|Brian        |         4973.0|        5028.0|2023-10-31 01:14:35|
47|Robert     |       1589.0|      1534.0|  55.0|Brian        |         7268.0|        7323.0|2023-10-31 01:14:35|
47|Robert     |       1589.0|      1534.0|  55.0|Brian        |         8926.0|        8981.0|2023-10-31 01:14:35|
47|Robert     |       7309.0|      7254.0|  55.0|Brian        |          709.0|         764.0|2023-10-31 01:14:35|
47|Robert     |       7309.0|      7254.0|  55.0|Brian        |         4973.0|        5028.0|2023-10-31 01:14:35|
47|Robert     |       7309.0|      7254.0|  55.0|Brian        |         7268.0|        7323.0|2023-10-31 01:14:35|
47|Robert     |       7309.0|      7254.0|  55.0|Brian        |         8926.0|        8981.0|2023-10-31 01:14:35|

12 row(s) fetched.
*/
