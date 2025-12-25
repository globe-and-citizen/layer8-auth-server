<template>
  <section class="bg-white rounded-2xl py-3 md:py-4 px-4 md:px-6 mb-6 md:mb-0 mt-6">
    <h1 class="font-medium text-lg md:text-xl text-black">Your usage statistics</h1>
    <p class="font-normal text-sm md:text-base text-[#8E8E93] mb-5">Your product data to
      use on your</p>

    <div class="grid grid-cols-4 gap-4">
      <div class="col-span-1 md:ml-20">
        <h2 class="text-xl font-bold">
          {{ stats.last_thirty_days_statistic?.total?.toFixed(6) }}
          {{ stats.unit_of_measurement }}
        </h2>
        <p class="text-gray-600 font-light">In Last 30 Days</p>
      </div>

      <div class="col-span-1 ml-10">
        <h2 class="text-xl font-bold">
          {{ stats.last_thirty_days_statistic?.average?.toFixed(6) }}
          {{ stats.unit_of_measurement }}
        </h2>
        <p class="s text-gray-600 font-light">Daily Average</p>
      </div>
      <div class="col-span-2">
        <div class="max-w-md mx-auto bg-gray-400 rounded-xl shadow-md overflow-hidden">
          <div class="p-4">
            <div class="uppercase tracking-wide text-lg font-extrabold text-gray-100">
              {{ stats.month_to_date?.month }}
            </div>
            <h3 class="text-white text-sm font-light">
              Month-to-date usage : {{
                stats.month_to_date?.month_to_date_usage?.toFixed(6)
              }} {{ stats.unit_of_measurement }}</h3>
            <h3 class="text-white text-sm font-light">
              Forecasted end-of-month usage : {{
                stats.month_to_date?.forecasted_end_of_month_usage?.toFixed(6)
              }} {{ stats.unit_of_measurement }}
            </h3>
          </div>
        </div>
      </div>
    </div>
    <div>
      <canvas id="statisticChart"></canvas>
      <p class="text-gray-600 font-ligh">The graph above presents statistics measured in
        kilobytes.</p>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from "vue";
import Chart from "chart.js/auto";

const stats = ref({
  metric_type: "",
  unit_of_measurement: "",
  month_to_date: {
    month: "",
    month_to_date_usage: 0,
    forecasted_end_of_month_usage: 0,
  },
  last_thirty_days_statistic: {
    total: 0.0,
    average: 0.0,
    details: [],
  },
})

onMounted(() => {
  // const statsResp = await fetch("/api/v1/usage-stats", {
  //   headers: { Authorization: `Bearer ${token.value}` },
  // })
  // stats.value = (await statsResp.json()).data

  stats.value = {
    metric_type: "data_transferred",
    unit_of_measurement: "GB",
    month_to_date: {
      month: "December",
      month_to_date_usage: 0,
      forecasted_end_of_month_usage: 0
    },
    last_thirty_days_statistic: {
      total: 0.007913035,
      average: 0.0002552591935483871,
      details: [
        {date: "Wed, 26 Nov 2025", total: 0},
        {date: "Thu, 27 Nov 2025", total: 0},
        {date: "Fri, 28 Nov 2025", total: 0.007913035},
        {date: "Sat, 29 Nov 2025", total: 0},
        {date: "Sun, 30 Nov 2025", total: 0},
        {date: "Mon, 01 Dec 2025", total: 0},
        {date: "Tue, 02 Dec 2025", total: 0},
        {date: "Wed, 03 Dec 2025", total: 0},
        {date: "Thu, 04 Dec 2025", total: 0},
        {date: "Fri, 05 Dec 2025", total: 0},
        {date: "Sat, 06 Dec 2025", total: 0},
        {date: "Sun, 07 Dec 2025", total: 0},
        {date: "Mon, 08 Dec 2025", total: 0},
        {date: "Tue, 09 Dec 2025", total: 0},
        {date: "Wed, 10 Dec 2025", total: 0},
        {date: "Thu, 11 Dec 2025", total: 0},
        {date: "Fri, 12 Dec 2025", total: 0},
        {date: "Sat, 13 Dec 2025", total: 0},
        {date: "Sun, 14 Dec 2025", total: 0},
        {
          date: "Mon, 15 Dec 2025",
          total: 0
        }, {date: "Tue, 16 Dec 2025", total: 0}, {
          date: "Wed, 17 Dec 2025",
          total: 0
        }, {date: "Thu, 18 Dec 2025", total: 0}, {
          date: "Fri, 19 Dec 2025",
          total: 0
        }, {date: "Sat, 20 Dec 2025", total: 0}, {
          date: "Sun, 21 Dec 2025",
          total: 0
        }, {date: "Mon, 22 Dec 2025", total: 0}, {
          date: "Tue, 23 Dec 2025",
          total: 0
        }, {date: "Wed, 24 Dec 2025", total: 0}, {
          date: "Thu, 25 Dec 2025",
          total: 0
        }, {date: "Thu, 25 Dec 2025", total: 0}]
    }
  }

  const mappedValue = stats.value.last_thirty_days_statistic?.details?.map(v => {
    v.total = (v.total * 1000000).toFixed(4)
    return v
  })

  const lineChart = new Chart(document.getElementById("statisticChart"), {
    type: 'line',
    data: {
      datasets: [
        {
          label: "Total Byte",
          borderWidth: 3,
          data: mappedValue,
        }
      ]
    },
    options: {
      parsing: {
        xAxisKey: 'date',
        yAxisKey: 'total'
      },
      plugins: {
        legend: false,
      },
    }
  })

});
</script>
