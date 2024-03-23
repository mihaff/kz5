function MoreChartOptions(){} 
var ChartData5 = {
	labels : ["Превосходно","Хорошо","Удовлетварительно","Есть над чем работать",],
	datasets : [
		{
		data : [496,455,320,48,],
		backgroundColor :['#a478f0','#37e881','#78c0f0','#A52A2A',],
		borderColor : ['rgba(136,136,136,0)','rgba(0,0,0,0)','rgba(156,89,182,0)','rgba(241,196,15,0)',],
		label:"Отзывы"},

]
	};
ChartOptions5= {
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
	text:'Оценки клиентов',
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
};
 DrawTheChart(ChartData5,ChartOptions5,"chart-05","pie");