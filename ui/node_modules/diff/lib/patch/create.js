/*istanbul ignore start*/'use strict';

exports.__esModule = true;
exports. /*istanbul ignore end*/structuredPatch = structuredPatch;
/*istanbul ignore start*/exports. /*istanbul ignore end*/createTwoFilesPatch = createTwoFilesPatch;
/*istanbul ignore start*/exports. /*istanbul ignore end*/createPatch = createPatch;

var /*istanbul ignore start*/_line = require('../diff/line') /*istanbul ignore end*/;

/*istanbul ignore start*/
function _toConsumableArray(arr) { if (Array.isArray(arr)) { for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) { arr2[i] = arr[i]; } return arr2; } else { return Array.from(arr); } }

/*istanbul ignore end*/function structuredPatch(oldFileName, newFileName, oldStr, newStr, oldHeader, newHeader, options) {
  if (!options) {
    options = { context: 4 };
  }

  var diff = /*istanbul ignore start*/(0, _line.diffLines) /*istanbul ignore end*/(oldStr, newStr);
  diff.push({ value: '', lines: [] }); // Append an empty value to make cleanup easier

  function contextLines(lines) {
    return lines.map(function (entry) {
      return ' ' + entry;
    });
  }

  var hunks = [];
  var oldRangeStart = 0,
      newRangeStart = 0,
      curRange = [],
      oldLine = 1,
      newLine = 1;
  /*istanbul ignore start*/
  var _loop = function _loop( /*istanbul ignore end*/i) {
    var current = diff[i],
        lines = current.lines || current.value.replace(/\n$/, '').split('\n');
    current.lines = lines;

    if (current.added || current.removed) {
      /*istanbul ignore start*/
      var _curRange;

      /*istanbul ignore end*/
      // If we have previous context, start with that
      if (!oldRangeStart) {
        var prev = diff[i - 1];
        oldRangeStart = oldLine;
        newRangeStart = newLine;

        if (prev) {
          curRange = options.context > 0 ? contextLines(prev.lines.slice(-options.context)) : [];
          oldRangeStart -= curRange.length;
          newRangeStart -= curRange.length;
        }
      }

      // Output our changes
      /*istanbul ignore start*/(_curRange = /*istanbul ignore end*/curRange).push. /*istanbul ignore start*/apply /*istanbul ignore end*/( /*istanbul ignore start*/_curRange /*istanbul ignore end*/, /*istanbul ignore start*/_toConsumableArray( /*istanbul ignore end*/lines.map(function (entry) {
        return (current.added ? '+' : '-') + entry;
      })));

      // Track the updated file position
      if (current.added) {
        newLine += lines.length;
      } else {
        oldLine += lines.length;
      }
    } else {
      // Identical context lines. Track line changes
      if (oldRangeStart) {
        // Close out any changes that have been output (or join overlapping)
        if (lines.length <= options.context * 2 && i < diff.length - 2) {
          /*istanbul ignore start*/
          var _curRange2;

          /*istanbul ignore end*/
          // Overlapping
          /*istanbul ignore start*/(_curRange2 = /*istanbul ignore end*/curRange).push. /*istanbul ignore start*/apply /*istanbul ignore end*/( /*istanbul ignore start*/_curRange2 /*istanbul ignore end*/, /*istanbul ignore start*/_toConsumableArray( /*istanbul ignore end*/contextLines(lines)));
        } else {
          /*istanbul ignore start*/
          var _curRange3;

          /*istanbul ignore end*/
          // end the range and output
          var contextSize = Math.min(lines.length, options.context);
          /*istanbul ignore start*/(_curRange3 = /*istanbul ignore end*/curRange).push. /*istanbul ignore start*/apply /*istanbul ignore end*/( /*istanbul ignore start*/_curRange3 /*istanbul ignore end*/, /*istanbul ignore start*/_toConsumableArray( /*istanbul ignore end*/contextLines(lines.slice(0, contextSize))));

          var hunk = {
            oldStart: oldRangeStart,
            oldLines: oldLine - oldRangeStart + contextSize,
            newStart: newRangeStart,
            newLines: newLine - newRangeStart + contextSize,
            lines: curRange
          };
          if (i >= diff.length - 2 && lines.length <= options.context) {
            // EOF is inside this hunk
            var oldEOFNewline = /\n$/.test(oldStr);
            var newEOFNewline = /\n$/.test(newStr);
            if (lines.length == 0 && !oldEOFNewline) {
              // special case: old has no eol and no trailing context; no-nl can end up before adds
              curRange.splice(hunk.oldLines, 0, '\\ No newline at end of file');
            } else if (!oldEOFNewline || !newEOFNewline) {
              curRange.push('\\ No newline at end of file');
            }
          }
          hunks.push(hunk);

          oldRangeStart = 0;
          newRangeStart = 0;
          curRange = [];
        }
      }
      oldLine += lines.length;
      newLine += lines.length;
    }
  };

  for (var i = 0; i < diff.length; i++) {
    /*istanbul ignore start*/
    _loop( /*istanbul ignore end*/i);
  }

  return {
    oldFileName: oldFileName, newFileName: newFileName,
    oldHeader: oldHeader, newHeader: newHeader,
    hunks: hunks
  };
}

function createTwoFilesPatch(oldFileName, newFileName, oldStr, newStr, oldHeader, newHeader, options) {
  var diff = structuredPatch(oldFileName, newFileName, oldStr, newStr, oldHeader, newHeader, options);

  var ret = [];
  if (oldFileName == newFileName) {
    ret.push('Index: ' + oldFileName);
  }
  ret.push('===================================================================');
  ret.push('--- ' + diff.oldFileName + (typeof diff.oldHeader === 'undefined' ? '' : '\t' + diff.oldHeader));
  ret.push('+++ ' + diff.newFileName + (typeof diff.newHeader === 'undefined' ? '' : '\t' + diff.newHeader));

  for (var i = 0; i < diff.hunks.length; i++) {
    var hunk = diff.hunks[i];
    ret.push('@@ -' + hunk.oldStart + ',' + hunk.oldLines + ' +' + hunk.newStart + ',' + hunk.newLines + ' @@');
    ret.push.apply(ret, hunk.lines);
  }

  return ret.join('\n') + '\n';
}

function createPatch(fileName, oldStr, newStr, oldHeader, newHeader, options) {
  return createTwoFilesPatch(fileName, fileName, oldStr, newStr, oldHeader, newHeader, options);
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIi4uLy4uL3NyYy9wYXRjaC9jcmVhdGUuanMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O2dDQUVnQixlLEdBQUEsZTt5REE4RkEsbUIsR0FBQSxtQjt5REF3QkEsVyxHQUFBLFc7O0FBeEhoQixJLHlCQUFBLCtCLHdCQUFBOzs7Ozt1QkFFTyxTQUFTLGVBQVQsQ0FBeUIsV0FBekIsRUFBc0MsV0FBdEMsRUFBbUQsTUFBbkQsRUFBMkQsTUFBM0QsRUFBbUUsU0FBbkUsRUFBOEUsU0FBOUUsRUFBeUYsT0FBekYsRUFBa0c7QUFDdkcsTUFBSSxDQUFDLE9BQUwsRUFBYztBQUNaLGNBQVUsRUFBRSxTQUFTLENBQVgsRUFBVjtBQUNEOztBQUVELE1BQU0sTyx5QkFBTyxvQix3QkFBQSxDQUFVLE1BQVYsRUFBa0IsTUFBbEIsQ0FBYjtBQUNBLE9BQUssSUFBTCxDQUFVLEVBQUMsT0FBTyxFQUFSLEVBQVksT0FBTyxFQUFuQixFQUFWLEU7O0FBRUEsV0FBUyxZQUFULENBQXNCLEtBQXRCLEVBQTZCO0FBQzNCLFdBQU8sTUFBTSxHQUFOLENBQVUsVUFBUyxLQUFULEVBQWdCO0FBQUUsYUFBTyxNQUFNLEtBQWI7QUFBcUIsS0FBakQsQ0FBUDtBQUNEOztBQUVELE1BQUksUUFBUSxFQUFaO0FBQ0EsTUFBSSxnQkFBZ0IsQ0FBcEI7QUFBQSxNQUF1QixnQkFBZ0IsQ0FBdkM7QUFBQSxNQUEwQyxXQUFXLEVBQXJEO0FBQUEsTUFDSSxVQUFVLENBRGQ7QUFBQSxNQUNpQixVQUFVLENBRDNCOztBQWJ1Ryw2Qix3QkFlOUYsQ0FmOEY7QUFnQnJHLFFBQU0sVUFBVSxLQUFLLENBQUwsQ0FBaEI7QUFBQSxRQUNNLFFBQVEsUUFBUSxLQUFSLElBQWlCLFFBQVEsS0FBUixDQUFjLE9BQWQsQ0FBc0IsS0FBdEIsRUFBNkIsRUFBN0IsRUFBaUMsS0FBakMsQ0FBdUMsSUFBdkMsQ0FEL0I7QUFFQSxZQUFRLEtBQVIsR0FBZ0IsS0FBaEI7O0FBRUEsUUFBSSxRQUFRLEtBQVIsSUFB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QUFLLEtBQXpCO0FBQ0Q7O0FBRUQsU0FBTyxJQUFJLElBQUosQ0FBUyxJQUFULElBQWlCLElBQXhCO0FBQ0Q7O0FBRU0sU0FBUyxXQUFULENBQXFCLFFBQXJCLEVBQStCLE1BQS9CLEVBQXVDLE1BQXZDLEVBQStDLFNBQS9DLEVBQTBELFNBQTFELEVBQXFFLE9BQXJFLEVBQThFO0FBQ25GLFNBQU8sb0JBQW9CLFFBQXBCLEVBQThCLFFBQTlCLEVBQXdDLE1BQXhDLEVBQWdELE1BQWhELEVBQXdELFNBQXhELEVBQW1FLFNBQW5FLEVBQThFLE9BQTlFLENBQVA7QUFDRCIsImZpbGUiOiJjcmVhdGUuanMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQge2RpZmZMaW5lc30gZnJvbSAnLi4vZGlmZi9saW5lJztcblxuZXhwb3J0IGZ1bmN0aW9uIHN0cnVjdHVyZWRQYXRjaChvbGRGaWxlTmFtZSwgbmV3RmlsZU5hbWUsIG9sZFN0ciwgbmV3U3RyLCBvbGRIZWFkZXIsIG5ld0hlYWRlciwgb3B0aW9ucykge1xuICBpZiAoIW9wdGlvbnMpIHtcbiAgICBvcHRpb25zID0geyBjb250ZXh0OiA0IH07XG4gIH1cblxuICBjb25zdCBkaWZmID0gZGlmZkxpbmVzKG9sZFN0ciwgbmV3U3RyKTtcbiAgZGlmZi5wdXNoKHt2YWx1ZTogJycsIGxpbmVzOiBbXX0pOyAgIC8vIEFwcGVuZCBhbiBlbXB0eSB2YWx1ZSB0byBtYWtlIGNsZWFudXAgZWFzaWVyXG5cbiAgZnVuY3Rpb24gY29udGV4dExpbmVzKGxpbmVzKSB7XG4gICAgcmV0dXJuIGxpbmVzLm1hcChmdW5jdGlvbihlbnRyeSkgeyByZXR1cm4gJyAnICsgZW50cnk7IH0pO1xuICB9XG5cbiAgbGV0IGh1bmtzID0gW107XG4gIGxldCBvbGRSYW5nZVN0YXJ0ID0gMCwgbmV3UmFuZ2VTdGFydCA9IDAsIGN1clJhbmdlID0gW10sXG4gICAgICBvbGRMaW5lID0gMSwgbmV3TGluZSA9IDE7XG4gIGZvciAobGV0IGkgPSAwOyBpIDwgZGlmZi5sZW5ndGg7IGkrKykge1xuICAgIGNvbnN0IGN1cnJlbnQgPSBkaWZmW2ldLFxuICAgICAgICAgIGxpbmVzID0gY3VycmVudC5saW5lcyB8fCBjdXJyZW50LnZhbHVlLnJlcGxhY2UoL1xcbiQvLCAnJykuc3BsaXQoJ1xcbicpO1xuICAgIGN1cnJlbnQubGluZXMgPSBsaW5lcztcblxuICAgIGlmIChjdXJyZW50LmFkZGVkIHx8IGN1cnJlbnQucmVtb3ZlZCkge1xuICAgICAgLy8gSWYgd2UgaGF2ZSBwcmV2aW91cyBjb250ZXh0LCBzdGFydCB3aXRoIHRoYXRcbiAgICAgIGlmICghb2xkUmFuZ2VTdGFydCkge1xuICAgICAgICBjb25zdCBwcmV2ID0gZGlmZltpIC0gMV07XG4gICAgICAgIG9sZFJhbmdlU3RhcnQgPSBvbGRMaW5lO1xuICAgICAgICBuZXdSYW5nZVN0YXJ0ID0gbmV3TGluZTtcblxuICAgICAgICBpZiAocHJldikge1xuICAgICAgICAgIGN1clJhbmdlID0gb3B0aW9ucy5jb250ZXh0ID4gMCA/IGNvbnRleHRMaW5lcyhwcmV2LmxpbmVzLnNsaWNlKC1vcHRpb25zLmNvbnRleHQpKSA6IFtdO1xuICAgICAgICAgIG9sZFJhbmdlU3RhcnQgLT0gY3VyUmFuZ2UubGVuZ3RoO1xuICAgICAgICAgIG5ld1JhbmdlU3RhcnQgLT0gY3VyUmFuZ2UubGVuZ3RoO1xuICAgICAgICB9XG4gICAgICB9XG5cbiAgICAgIC8vIE91dHB1dCBvdXIgY2hhbmdlc1xuICAgICAgY3VyUmFuZ2UucHVzaCguLi4gbGluZXMubWFwKGZ1bmN0aW9uKGVudHJ5KSB7XG4gICAgICAgIHJldHVybiAoY3VycmVudC5hZGRlZCA/ICcrJyA6ICctJykgKyBlbnRyeTtcbiAgICAgIH0pKTtcblxuICAgICAgLy8gVHJhY2sgdGhlIHVwZGF0ZWQgZmlsZSBwb3NpdGlvblxuICAgICAgaWYgKGN1cnJlbnQuYWRkZWQpIHtcbiAgICAgICAgbmV3TGluZSArPSBsaW5lcy5sZW5ndGg7XG4gICAgICB9IGVsc2Uge1xuICAgICAgICBvbGRMaW5lICs9IGxpbmVzLmxlbmd0aDtcbiAgICAgIH1cbiAgICB9IGVsc2Uge1xuICAgICAgLy8gSWRlbnRpY2FsIGNvbnRleHQgbGluZXMuIFRyYWNrIGxpbmUgY2hhbmdlc1xuICAgICAgaWYgKG9sZFJhbmdlU3RhcnQpIHtcbiAgICAgICAgLy8gQ2xvc2Ugb3V0IGFueSBjaGFuZ2VzIHRoYXQgaGF2ZSBiZWVuIG91dHB1dCAob3Igam9pbiBvdmVybGFwcGluZylcbiAgICAgICAgaWYgKGxpbmVzLmxlbmd0aCA8PSBvcHRpb25zLmNvbnRleHQgKiAyICYmIGkgPCBkaWZmLmxlbmd0aCAtIDIpIHtcbiAgICAgICAgICAvLyBPdmVybGFwcGluZ1xuICAgICAgICAgIGN1clJhbmdlLnB1c2goLi4uIGNvbnRleHRMaW5lcyhsaW5lcykpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgIC8vIGVuZCB0aGUgcmFuZ2UgYW5kIG91dHB1dFxuICAgICAgICAgIGxldCBjb250ZXh0U2l6ZSA9IE1hdGgubWluKGxpbmVzLmxlbmd0aCwgb3B0aW9ucy5jb250ZXh0KTtcbiAgICAgICAgICBjdXJSYW5nZS5wdXNoKC4uLiBjb250ZXh0TGluZXMobGluZXMuc2xpY2UoMCwgY29udGV4dFNpemUpKSk7XG5cbiAgICAgICAgICBsZXQgaHVuayA9IHtcbiAgICAgICAgICAgIG9sZFN0YXJ0OiBvbGRSYW5nZVN0YXJ0LFxuICAgICAgICAgICAgb2xkTGluZXM6IChvbGRMaW5lIC0gb2xkUmFuZ2VTdGFydCArIGNvbnRleHRTaXplKSxcbiAgICAgICAgICAgIG5ld1N0YXJ0OiBuZXdSYW5nZVN0YXJ0LFxuICAgICAgICAgICAgbmV3TGluZXM6IChuZXdMaW5lIC0gbmV3UmFuZ2VTdGFydCArIGNvbnRleHRTaXplKSxcbiAgICAgICAgICAgIGxpbmVzOiBjdXJSYW5nZVxuICAgICAgICAgIH07XG4gICAgICAgICAgaWYgKGkgPj0gZGlmZi5sZW5ndGggLSAyICYmIGxpbmVzLmxlbmd0aCA8PSBvcHRpb25zLmNvbnRleHQpIHtcbiAgICAgICAgICAgIC8vIEVPRiBpcyBpbnNpZGUgdGhpcyBodW5rXG4gICAgICAgICAgICBsZXQgb2xkRU9GTmV3bGluZSA9ICgvXFxuJC8udGVzdChvbGRTdHIpKTtcbiAgICAgICAgICAgIGxldCBuZXdFT0ZOZXdsaW5lID0gKC9cXG4kLy50ZXN0KG5ld1N0cikpO1xuICAgICAgICAgICAgaWYgKGxpbmVzLmxlbmd0aCA9PSAwICYmICFvbGRFT0ZOZXdsaW5lKSB7XG4gICAgICAgICAgICAgIC8vIHNwZWNpYWwgY2FzZTogb2xkIGhhcyBubyBlb2wgYW5kIG5vIHRyYWlsaW5nIGNvbnRleHQ7IG5vLW5sIGNhbiBlbmQgdXAgYmVmb3JlIGFkZHNcbiAgICAgICAgICAgICAgY3VyUmFuZ2Uuc3BsaWNlKGh1bmsub2xkTGluZXMsIDAsICdcXFxcIE5vIG5ld2xpbmUgYXQgZW5kIG9mIGZpbGUnKTtcbiAgICAgICAgICAgIH0gZWxzZSBpZiAoIW9sZEVPRk5ld2xpbmUgfHwgIW5ld0VPRk5ld2xpbmUpIHtcbiAgICAgICAgICAgICAgY3VyUmFuZ2UucHVzaCgnXFxcXCBObyBuZXdsaW5lIGF0IGVuZCBvZiBmaWxlJyk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgfVxuICAgICAgICAgIGh1bmtzLnB1c2goaHVuayk7XG5cbiAgICAgICAgICBvbGRSYW5nZVN0YXJ0ID0gMDtcbiAgICAgICAgICBuZXdSYW5nZVN0YXJ0ID0gMDtcbiAgICAgICAgICBjdXJSYW5nZSA9IFtdO1xuICAgICAgICB9XG4gICAgICB9XG4gICAgICBvbGRMaW5lICs9IGxpbmVzLmxlbmd0aDtcbiAgICAgIG5ld0xpbmUgKz0gbGluZXMubGVuZ3RoO1xuICAgIH1cbiAgfVxuXG4gIHJldHVybiB7XG4gICAgb2xkRmlsZU5hbWU6IG9sZEZpbGVOYW1lLCBuZXdGaWxlTmFtZTogbmV3RmlsZU5hbWUsXG4gICAgb2xkSGVhZGVyOiBvbGRIZWFkZXIsIG5ld0hlYWRlcjogbmV3SGVhZGVyLFxuICAgIGh1bmtzOiBodW5rc1xuICB9O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gY3JlYXRlVHdvRmlsZXNQYXRjaChvbGRGaWxlTmFtZSwgbmV3RmlsZU5hbWUsIG9sZFN0ciwgbmV3U3RyLCBvbGRIZWFkZXIsIG5ld0hlYWRlciwgb3B0aW9ucykge1xuICBjb25zdCBkaWZmID0gc3RydWN0dXJlZFBhdGNoKG9sZEZpbGVOYW1lLCBuZXdGaWxlTmFtZSwgb2xkU3RyLCBuZXdTdHIsIG9sZEhlYWRlciwgbmV3SGVhZGVyLCBvcHRpb25zKTtcblxuICBjb25zdCByZXQgPSBbXTtcbiAgaWYgKG9sZEZpbGVOYW1lID09IG5ld0ZpbGVOYW1lKSB7XG4gICAgcmV0LnB1c2goJ0luZGV4OiAnICsgb2xkRmlsZU5hbWUpO1xuICB9XG4gIHJldC5wdXNoKCc9PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09Jyk7XG4gIHJldC5wdXNoKCctLS0gJyArIGRpZmYub2xkRmlsZU5hbWUgKyAodHlwZW9mIGRpZmYub2xkSGVhZGVyID09PSAndW5kZWZpbmVkJyA/ICcnIDogJ1xcdCcgKyBkaWZmLm9sZEhlYWRlcikpO1xuICByZXQucHVzaCgnKysrICcgKyBkaWZmLm5ld0ZpbGVOYW1lICsgKHR5cGVvZiBkaWZmLm5ld0hlYWRlciA9PT0gJ3VuZGVmaW5lZCcgPyAnJyA6ICdcXHQnICsgZGlmZi5uZXdIZWFkZXIpKTtcblxuICBmb3IgKGxldCBpID0gMDsgaSA8IGRpZmYuaHVua3MubGVuZ3RoOyBpKyspIHtcbiAgICBjb25zdCBodW5rID0gZGlmZi5odW5rc1tpXTtcbiAgICByZXQucHVzaChcbiAgICAgICdAQCAtJyArIGh1bmsub2xkU3RhcnQgKyAnLCcgKyBodW5rLm9sZExpbmVzXG4gICAgICArICcgKycgKyBodW5rLm5ld1N0YXJ0ICsgJywnICsgaHVuay5uZXdMaW5lc1xuICAgICAgKyAnIEBAJ1xuICAgICk7XG4gICAgcmV0LnB1c2guYXBwbHkocmV0LCBodW5rLmxpbmVzKTtcbiAgfVxuXG4gIHJldHVybiByZXQuam9pbignXFxuJykgKyAnXFxuJztcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNyZWF0ZVBhdGNoKGZpbGVOYW1lLCBvbGRTdHIsIG5ld1N0ciwgb2xkSGVhZGVyLCBuZXdIZWFkZXIsIG9wdGlvbnMpIHtcbiAgcmV0dXJuIGNyZWF0ZVR3b0ZpbGVzUGF0Y2goZmlsZU5hbWUsIGZpbGVOYW1lLCBvbGRTdHIsIG5ld1N0ciwgb2xkSGVhZGVyLCBuZXdIZWFkZXIsIG9wdGlvbnMpO1xufVxuIl19
