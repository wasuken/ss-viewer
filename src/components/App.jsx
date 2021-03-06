import React, { Component } from 'react';

class App extends Component{
    style = { width: "1.2em",
              lineHeight: "1.2em",
              listStyleType: "none",
              border: "1px solid blue",
              float: "left",
              margin: "0.1em",
              padding: "0px",
              fontFamily: "Arial,sans-serif",
              fontWeight: "bold",
              textAlign: "center"}
    constructor(props){
        super(props);
        this.state = {
            json: [],
            displayContents: "",
            page: 0
        };
        this.getJson();
        this.clickAction = this.clickAction.bind(this);
        let that = this;
        fetch("/api/v1/ss/page").then(resp=>resp.json()).then(json =>{
            that.setState({
                json: [],
                displayContents: "",
                page: json.page
            });
        });
    }
    clickAction(){
        let result = this.state.json.find((map)=>{
            return map.Id == document.querySelector("#select").value;
        });
        this.setState({
            json: this.state.json,
            displayContents: result.Contents,
            page: this.state.page
        });
    }
    getJson(page=1){
        let that = this;
        fetch("/api/v1/ss?page=" + page).then((resp)=>resp.json()).then((json)=>{
            that.setState({
                json: json,
                displayContents: this.state.contents,
                page: this.state.page
            });
        });
    }
    getLink(){
        let lst = [];
        let page = this.state.page > 10? 10: this.state.page;
        for(let i=1;i<page+1;i++){
            lst.push(<a className="page" href="#" onClick={() => this.getJson(i)}><li style={this.style}>{i}</li></a>);
        }
        return lst;
    }
    render(){
        let lst = this.getLink();
        
        let titles = this.state.json.map((map)=>(
            <option value={map.Id}>{map.Title}</option>
        ));
        return (
            <div>
              <div style={{float:"left", width:"30%"}}>
                <select id="select" onChange={this.clickAction} size="10">
                  {titles}
                </select>
                <div>
                  {lst}
                </div>
              </div>
              <div style={{float:"left", width:"70%", height: "100%",overflowY: "scroll"}}
                   className="text" dangerouslySetInnerHTML={{__html: this.state.displayContents}}>
              </div>
            </div>
        );
    }
}

export default App;
