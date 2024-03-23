function MoreChartOptions(){} 
var ChartData4 = {
	labels : ["Classification","Regression",],
	datasets : [
		{
		data : [389,455],
		backgroundColor :['#a478f0','#78c0f0'],
		borderColor : ['rgba(136,136,136,0)','rgba(0,0,0,0)'],
		label:"Стартап"},

]
	};
ChartOptions4= {
responsive:false,
	layout:{padding:{top:12,left:12,bottom:12,},},
	plugins:{
datalabels:{display:true,
	color:'rgba(0,0,0,0)',
	font:{
		style:' bold',},},
},
legend:{
	position:'bottom',
	labels:{
		usePointStyle:true,

		generateLabels: function(chart){
			 return  chart.data.labels.map( function( dataset, i ){
				return{
					text:dataset,
					lineCap:chart.data.datasets[0].borderCapStyle,
					lineDash:[],
					lineDashOffset: 0,
					lineJoin:chart.data.datasets[0].borderJoinStyle,
					pointStyle:'circle',
					fillStyle:chart.data.datasets[0].backgroundColor[i],
					strokeStyle:chart.data.datasets[0].borderColor[i],
					lineWidth:chart.data.datasets[0].pointBorderWidth,
					lineDash:chart.data.datasets[0].borderDash,
				}
			})
		},

	},
},

title:{
	display:true,
	text:'Использование моделей за 2023 г.',
	fontColor:'#3498db',
	fontSize:32,
	fontStyle:' bold',
	},
elements: {
	arc: {borderWidth:0.01,
},
	line: {
},
	rectangle: {
},
},
tooltips:{
},
hover:{
	mode:'nearest',
	animationDuration:400,
},
animation:{
	easing:'easeInQuart',
	},
};
 DrawTheChart(ChartData4,ChartOptions4,"chart-04","doughnut");