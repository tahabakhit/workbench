---
name: explain-gap
description: Use when two numbers that should match do not - a dashboard versus its source table, two reports, a metric before versus after a change, Snowflake versus Tableau. Finds the cause of the gap with queries. Not for failing code or builds; that is diagnose.
---

# Explain gap

## 1. Pin

Write down for both numbers: the exact value, where it comes from (table,
view, dashboard, query), date range, timezone, filters and grain (one row per
what). Find the metric definition before assuming one, in the semantic layer, metric
catalog or dbt docs the instructions or tools point to.

Done when both sides are pinned and any difference in the pins is listed.

## 2. Reproduce the gap

Write one query that returns both numbers side by side for the pinned window.
Nothing else starts until it runs and shows the gap.

Done when the query output shows the gap.

## 3. Check known causes, cheapest first

| Cause | Test |
|---|---|
| Timezone or date boundary | Shift the window by the UTC offset; does the gap move? |
| Join fan-out | Compare `count(*)` with `count(distinct <key>)` after each join |
| Late or reprocessed data | Rerun on a window that closed a week ago |
| Filter difference | Remove filters one at a time on the larger side |
| Definition drift | Compare the two SQL definitions line by line |
| Duplicates or test data | Group by key having `count(*) > 1`; check test flags |

## 4. Localise

Split the gap by one dimension at a time (day, market, category, platform)
until the gap sits in one slice.

## 5. Report

State the cause, the query that proves it, how much of the gap it explains,
and any rest still unexplained.

Done when the explained part is shown by query output, not reasoning.
