const combobox = ({
    ajax,
    searchCallback,
    onChangeCallback,
    id = 'id',
    text = 'text',
    data,
    isMultiple = false,
    isInsertable = false,
    isDisplayItem = false,
    displaySearch = true,
    label = null,
    searchTag = null,
    searchPlaceholder = null
}) => {
    return { 
        // STATES
        options : [
            // { id : 0, text : 'xxxxx'}
        ],
        data : [
            // { id : 0, text : 'xxxxx'}
        ],
        selected : [
            // { id : 0, text : 'xxxxx'}
        ],
        isMultiple : false,
        isShowOptions : false,
        isInsertable : false,
        isDisplayItem : false,
        displaySearch : true,
        searchTag : null,
        searchPlaceholder : null,
        label : '',
        // INITIALIZED
        async init() {
            this.label = label;
            this.isInsertable = isInsertable;
            this.isDisplayItem = isDisplayItem;
            this.displaySearch = displaySearch;
            this.searchTag = searchTag;
            this.searchPlaceholder = searchPlaceholder;

            // SET DATA
            if(data != undefined){
                this.data = data.map((i) => {
                    return {id : i[id], text : i[text]}
                })
            }
            else if(ajax != undefined){
                const items = await ajax()
                this.data = items.map((i) => {
                    return {id : i[id], text : i[text]}
                })
            }
            else {
                console.log("INFO : NO INITIAL DATA SET")
            }
            this.options = this.data
        },
        // EVENT HANDLER
        onInputHandler(e) {
            this.setOptions(e.target.value)
        },
        onFocusIn() {
            this.isShowOptions = true
        },
        onClickOption(){
            this.isShowOptions = !this.isShowOptions
        },
        onInsertItem(e){
            if (!this.isInsertable)
                return

            const value = e.target.value;

            if(!this.selected.find(e => e.text === value)) {
                this.insertSelectedItem({ id : 0, text : value})
            }

            e.target.value = '';
            this.options = this.data;
        },
        onSelectOption(item){
            if(onChangeCallback != undefined){
                onChangeCallback(item)
            }

            if(this.selected.find(v => v.id === item.id)) {
                this.removeSelectedItem(item)
                return;
            }
            this.insertSelectedItem(item)
        },
        onUnselectOption(item){
            this.removeSelectedItem(item)
        },
        isSelected(id) {
            return this.selected.some(v => v.id === id)
        },
        // METHODS
        async setOptions(value){
            if(searchCallback != undefined) {
                const result = await searchCallback({search:value})
                
                if (result == null){
                    this.options = []
                    return
                }

                this.options = result.map((i) => {
                    return {id : i[id], text : i[text]}
                })
                return
            }

            this.options = this.data.filter((v, i) => { return v.text.toLowerCase().includes(value.toLowerCase())})
        },
        insertSelectedItem(item){
            if(!isMultiple)
                this.selected = []
            
            this.selected.push(item)
        },
        removeSelectedItem({id, text}){
            this.selected = this.isMultiple ? [] : this.selected.filter(v => (v.id != id || v.text != text))
        },
        template : `<div @click.outside="isShowOptions=false">
                        <template x-if="label != null">
                            <label class="block text-sm font-medium text-gray-700" x-text="label"></label>
                        </template>
                        <div class="relative mt-1">
                            <input x-model="selected.map(v => v.text).join()" type="text" class="w-full rounded-md border border-gray-300 bg-white py-2 pl-3 pr-12 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm" role="combobox" aria-controls="options" aria-expanded="false" x-on:focusin="onFocusIn" readonly>
                            <button x-on:click="onClickOption" type="button" class="absolute inset-y-0 right-0 flex items-center rounded-r-md px-2 focus:outline-none">
                                <!-- Heroicon name: solid/selector -->
                                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                                    <path fill-rule="evenodd" d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                                </svg>
                            </button>
                            <ul x-show="isShowOptions" class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm" id="options" role="listbox">
                            <template x-if='displaySearch'>
                                <li class="p-3">
                                <input x-ref="filter" type="text" :placeholder="searchPlaceholder" class="w-full rounded-md border border-gray-300 bg-white py-2 pl-3 pr-12 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm" role="combobox" aria-controls="options" aria-expanded="false" @input.debounce.1000ms="onInputHandler" @keyup.enter='onInsertItem'>
                                    <template x-if="searchTag != null">
                                        <small class="text-gray-700" x-text="searchTag"></small>
                                    </template>
                                </li>
                            </template>
                            <template x-for='item in options' :key="item.id">
                                <!--
                                Combobox option, manage highlight styles based on mouseenter/mouseleave and keyboard navigation.
                                Active: "text-white bg-indigo-600", Not Active: "text-gray-900"
                                -->
                                <li :id="item.id" :value="item.id" x-on:click="(e) => { onSelectOption(item) }" class="relative cursor-default select-none border-b py-2 pl-3 pr-9 text-gray-900 hover:bg-gray-100" role="option" tabindex="-1">
                                    <!-- Selected: "font-semibold" -->
                                    <span class="block truncate" x-text="item.text"></span>
                                    <!--
                                    Checkmark, only display for selected option.
                                    Active: "text-white", Not Active: "text-indigo-600"
                                    -->
                                    <span x-show="isSelected(item.id)" class="absolute inset-y-0 right-0 flex items-center pr-4 text-indigo-600">
                                        <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    </span>
                                </li>
                            </template>
                            </ul>
                        </div>
                        <div class='mt-1 bg-gray-50 p-1 rounded-md' x-show="(selected.length > 0) && isDisplayItem">
                            <template x-for='item in selected'>
                                <span class="inline-flex items-center py-0.5 pl-2 pr-0.5 rounded-full text-xs font-medium" :class="item.id == 0 ? 'bg-red-200 text-red-700' : 'bg-gray-200 text-gray-700'">
                                    <span x-text="item.text"></span>
                                    <button type="button" class="flex-shrink-0 ml-0.5 h-4 w-4 rounded-full inline-flex items-center justify-center focus:outline-none focus:text-white" :class="item.id == 0 ? 'text-red-400 hover:bg-red-200 hover:text-red-500 focus:bg-red-500' : 'text-gray-400 hover:bg-gray-200 hover:text-gray-500 focus:bg-gray-500'" @click="onUnselectOption(item)">
                                        <span class="sr-only">Remove small option</span>
                                        <svg class="h-2 w-2" stroke="currentColor" fill="none" viewBox="0 0 8 8">
                                            <path stroke-linecap="round" stroke-width="1.5" d="M1 1l6 6m0-6L1 7" />
                                        </svg>
                                    </button>
                                </span>
                            </template>
                        </<div>
                    </div>`
    }
}